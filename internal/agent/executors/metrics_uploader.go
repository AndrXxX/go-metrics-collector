package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

type metricsUploader struct {
	rs *requestsender.RequestSender
}

func (c *metricsUploader) Execute(result dto.MetricsDto) error {
	for metric, value := range result.Gauge {
		params := types.URLParams{"metricType": metrics.Gauge, "metric": metric, "value": value}
		_ = c.rs.Post(params, "text/plain")
	}
	for metric, value := range result.Counter {
		params := types.URLParams{"metricType": metrics.Counter, "metric": metric, "value": value}
		_ = c.rs.Post(params, "text/plain")
	}
	return nil
}

func NewUploader(rs *requestsender.RequestSender) *metricsUploader {
	return &metricsUploader{rs}
}
