package metricsuploader

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

type plainTextMetricsUploader struct {
	rs *requestsender.RequestSender
	ub urlBuilder
}

func (c *plainTextMetricsUploader) Execute(result dto.MetricsDto) error {
	for metric, value := range result.Gauge {
		params := types.URLParams{"metricType": metrics.Gauge, "metric": metric, "value": value}
		url := c.ub.Build(params)
		_ = c.rs.Post(url, contenttypes.TextPlain, nil)
	}
	for metric, value := range result.Counter {
		params := types.URLParams{"metricType": metrics.Counter, "metric": metric, "value": value}
		url := c.ub.Build(params)
		_ = c.rs.Post(url, contenttypes.TextPlain, nil)
	}
	return nil
}

func NewPlainTextUploader(rs *requestsender.RequestSender, ub urlBuilder) *plainTextMetricsUploader {
	return &plainTextMetricsUploader{rs, ub}
}
