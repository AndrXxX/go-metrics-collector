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

func (c *plainTextMetricsUploader) execute(result dto.MetricsDto) error {
	for _, metric := range result.All() {
		var value any
		if metric.MType == metrics.Gauge {
			value = *metric.Value
		} else {
			value = *metric.Delta
		}
		params := types.URLParams{"metricType": metric.MType, "metric": metric.ID, "value": value}
		url := c.ub.Build(params)
		_ = c.rs.Post(url, contenttypes.TextPlain, nil)
	}
	return nil
}

// Process выполняет загрузку метрик
func (c *plainTextMetricsUploader) Process(results <-chan dto.MetricsDto) error {
	for result := range results {
		err := c.execute(result)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewPlainTextUploader возвращает сервис plainTextMetricsUploader для загрузки метрик с помощью урла
func NewPlainTextUploader(rs *requestsender.RequestSender, ub urlBuilder) *plainTextMetricsUploader {
	return &plainTextMetricsUploader{rs, ub}
}
