package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/hashgenerator"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricscollector"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/scheduler"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Run(config *config.Config) error {
	m := dto.NewMetricsDto()
	s := scheduler.NewIntervalScheduler(config.Intervals.SleepInterval)
	s.Add(metricscollector.New(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)

	ub := metricurlbuilder.New(config.Common.Host)
	hg, err := hashgenerator.New(config.Common.Key)
	if err != nil {
		logger.Log.Error("Error on create hash generator", zap.Error(err))
	}
	rs := requestsender.New(http.DefaultClient, hg)
	s.Add(metricsuploader.NewJSONUploader(rs, ub, config.Intervals.RepeatIntervals), time.Duration(config.Intervals.ReportInterval)*time.Second)

	return s.Run(*m)
}
