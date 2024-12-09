package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

const shutdownTimeout = 5 * time.Second

type agent struct {
	c          *config.Config
	collectors types.ItemsList[scheduler.Collector]
	processors types.ItemsList[scheduler.Processor]
}

func New(c *config.Config, opts ...Option) *agent {
	a := &agent{
		c:          c,
		collectors: make(types.ItemsList[scheduler.Collector], 0),
		processors: make(types.ItemsList[scheduler.Processor], 0),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Run запускает агента
func (a *agent) Run(commonCtx context.Context) error {
	ctx, cancel := context.WithCancel(commonCtx)
	defer cancel()
	s := scheduler.NewIntervalScheduler(time.Duration(a.c.Intervals.SleepInterval) * time.Second)

	for _, collector := range a.collectors {
		s.AddCollector(collector, time.Duration(a.c.Intervals.PollInterval)*time.Second)
	}
	for _, processor := range a.processors {
		for count := a.c.Common.RateLimit; count > 0; count-- {
			s.AddProcessor(processor, time.Duration(a.c.Intervals.ReportInterval)*time.Second)
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
