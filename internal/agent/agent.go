package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/client"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/compressor"
	grpsserv "github.com/AndrXxX/go-metrics-collector/internal/agent/services/grpc"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/grpc/dealoptions"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/runtimemetricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/tlsconfig"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/vmmetricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

const shutdownTimeout = 5 * time.Second

// Run запускает агента
func Run(commonCtx context.Context, config *config.Config) error {
	ctx, cancel := context.WithCancel(commonCtx)
	defer cancel()
	s := scheduler.NewIntervalScheduler(time.Duration(config.Intervals.SleepInterval) * time.Second)

	for _, collector := range getCollectors(config) {
		s.AddCollector(collector, time.Duration(config.Intervals.PollInterval)*time.Second)
	}

	if processors, err := getProcessors(config); err != nil {
		return err
	} else {
		for _, processor := range processors {
			for count := config.Common.RateLimit; count > 0; count-- {
				s.AddProcessor(processor, time.Duration(config.Intervals.ReportInterval)*time.Second)
			}
		}
	}

	go func() {
		err := s.Run()
		if err != nil {
			logger.Log.Error("Failed to run scheduler", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Log.Info("shutting down agent gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdown := make(chan struct{}, 1)
	go func() {
		err := s.Shutdown(shutdownCtx)
		if err != nil {
			logger.Log.Error("Failed to shutdown scheduler", zap.Error(err))
		}
		shutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("agent shutdown: %w", shutdownCtx.Err())
	case <-shutdown:
		log.Println("finished")
	}

	return nil
}

func getCollectors(config *config.Config) []scheduler.Collector {
	var list []scheduler.Collector
	list = append(list, runtimemetricscollector.New(&config.Metrics))
	list = append(list, vmmetricscollector.New())
	return list
}

func getProcessors(config *config.Config) ([]scheduler.Processor, error) {
	hg := hashgenerator.Factory().SHA256()
	var list []scheduler.Processor
	if config.Common.GRPCHost != "" {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
		opts = append(opts, grpc.WithUnaryInterceptor(dealoptions.WithXRealIP(config.Common.Host)))
		opts = append(opts, grpc.WithUnaryInterceptor(dealoptions.WithSHA256(hg, config.Common.Key)))
		updater := grpsserv.NewGRPCMetricsUpdater(config.Common.GRPCHost, opts)
		list = append(list, metricsuploader.NewGRPCUploader(updater))
	}
	if config.Common.Host != "" && len(list) == 0 {
		ub := metricurlbuilder.New(config.Common.Host)

		httpClient, err := client.Provider{ConfProvider: tlsconfig.Provider{CryptoKeyPath: config.Common.CryptoKey}}.Fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch client: %w", err)
		}
		rs := requestsender.New(
			httpClient,
			requestsender.WithGzip(compressor.GzipCompressor{}),
			requestsender.WithSHA256(hg, config.Common.Key),
			requestsender.WithXRealIP(config.Common.Host),
		)
		processor := metricsuploader.NewJSONUploader(rs, ub, config.Intervals.RepeatIntervals)
		list = append(list, processor)
	}
	return list, nil
}
