package agent

import (
	"context"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/runtimemetricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/vmmetricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

const shutdownTimeout = 5 * time.Second

func Run(commonCtx context.Context, config *config.Config) error {
	ctx, cancel := context.WithCancel(commonCtx)
	defer cancel()
	s := scheduler.NewIntervalScheduler(time.Duration(config.Intervals.SleepInterval) * time.Second)
	rmc := runtimemetricscollector.New(&config.Metrics)
	s.AddCollector(rmc, time.Duration(config.Intervals.PollInterval)*time.Second)

	vmc := vmmetricscollector.New()
	s.AddCollector(vmc, time.Duration(config.Intervals.PollInterval)*time.Second)

	ub := metricurlbuilder.New(config.Common.Host)
	hg := hashgenerator.Factory().SHA256()
	rs := requestsender.New(http.DefaultClient, hg, config.Common.Key)
	for count := config.Common.RateLimit; count > 0; count-- {
		processor := metricsuploader.NewJSONUploader(rs, ub, config.Intervals.RepeatIntervals)
		s.AddProcessor(processor, time.Duration(config.Intervals.ReportInterval)*time.Second)
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
