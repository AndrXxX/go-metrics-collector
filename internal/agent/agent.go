package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/executors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsсollector"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"net/http"
	"time"
)

func Run(config *config.Config) error {
	m := dto.NewMetricsDto()
	s := scheduler.NewIntervalScheduler(config.Intervals.SleepInterval)
	s.Add(metricsсollector.New(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)

	rs := requestsender.New(metricurlbuilder.New(config.Common.Host), http.DefaultClient)
	s.Add(executors.NewUploader(rs), time.Duration(config.Intervals.ReportInterval)*time.Second)

	err := s.Run(*m)
	return err
}
