package metricsuploader

import (
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
)

type JSONMetrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

type jsonMetricsUploader struct {
	rs *requestsender.RequestSender
	ub urlBuilder
}

func (c *jsonMetricsUploader) Execute(result dto.MetricsDto) error {
	for metric, value := range result.Gauge {
		err := c.send(JSONMetrics{
			ID:    metric,
			MType: metrics.Gauge,
			Value: &value,
		})
		if err != nil {
			logger.Log.Error("error send response", zap.Error(err))
		}
	}
	for metric, value := range result.Counter {
		err := c.send(JSONMetrics{
			ID:    metric,
			MType: metrics.Counter,
			Delta: &value,
		})
		if err != nil {
			logger.Log.Error("error send response", zap.Error(err))
		}
	}
	return nil
}

func (c *jsonMetricsUploader) send(m JSONMetrics) error {
	url := c.ub.Build(types.URLParams{})
	encoded, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return c.rs.Post(url, contenttypes.ApplicationJSON, encoded)
}

func NewJSONUploader(rs *requestsender.RequestSender, ub urlBuilder) *jsonMetricsUploader {
	return &jsonMetricsUploader{rs, ub}
}
