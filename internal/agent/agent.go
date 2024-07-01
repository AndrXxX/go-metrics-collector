package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"net/http"
	"time"
)

func Run(config *config.Config) error {
	m := dto.NewMetricsDto()
	s := scheduler.NewIntervalScheduler(config.Intervals.SleepInterval)
	s.Add(metricscollector.New(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)

	ub := metricurlbuilder.New(config.Common.Host)
	rs := requestsender.New(http.DefaultClient)
	s.Add(metricsuploader.NewJsonUploader(rs, ub), time.Duration(config.Intervals.ReportInterval)*time.Second)

	err := s.Run(*m)
	return err
}
