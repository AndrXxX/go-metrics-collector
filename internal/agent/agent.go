package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/executors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"time"
)

func Run(config *config.Config) error {
	m := metrics.NewMetrics()
	scheduler := executors.NewIntervalScheduler()
	scheduler.Add(executors.NewCollector(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)
	scheduler.Add(executors.NewUploader(config.Common.Host), time.Duration(config.Intervals.ReportInterval)*time.Second)
	err := scheduler.Run(*m)
	return err
}
