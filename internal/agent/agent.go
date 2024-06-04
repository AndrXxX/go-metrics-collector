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
	scheduler.Add(executors.NewCollector(time.Duration(config.Intervals.PollInterval), &config.Metrics), time.Duration(config.Intervals.PollInterval))
	scheduler.Add(executors.NewUploader(time.Duration(config.Intervals.ReportInterval), &config.Common.Host), time.Duration(config.Intervals.ReportInterval))
	err := scheduler.Run(*m)
	return err
}
