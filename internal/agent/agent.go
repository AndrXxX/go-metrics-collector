package agent

import (
	"context"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"log"
	"net/http"
	"time"
)

const shutdownTimeout = 5 * time.Second

func Run(commonCtx context.Context, config *config.Config) error {
	m := dto.NewMetricsDto()
	s := scheduler.NewIntervalScheduler(config.Intervals.SleepInterval)
	s.Add(metricscollector.New(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)

	ub := metricurlbuilder.New(config.Common.Host)
	hg := hashgenerator.New()
	rs := requestsender.New(http.DefaultClient, hg, config.Common.Key)
	s.Add(metricsuploader.NewJSONUploader(rs, ub, config.Intervals.RepeatIntervals), time.Duration(config.Intervals.ReportInterval)*time.Second)

	err := s.Run(*m)
	if err != nil {
		return err
	}

	<-commonCtx.Done()
	logger.Log.Info("shutting down agent gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdown := make(chan struct{}, 1)
	go func() {
		// TODO: stop scheduler with shutdownCtx
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
