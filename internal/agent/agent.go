package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/executors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/utils"
	"net/http"
	"time"
)

func Run(config *config.Config) error {
	m := metrics.NewMetrics()
	s := scheduler.NewIntervalScheduler(config.Intervals.SleepInterval)
	s.Add(executors.NewCollector(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)

	rs := utils.NewRequestSender(utils.NewMetricURLBuilder(config.Common.Host), http.DefaultClient)
	s.Add(executors.NewUploader(rs), time.Duration(config.Intervals.ReportInterval)*time.Second)

	err := s.Run(*m)
	return err
}
