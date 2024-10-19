package metricsuploader

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type jsonMetricsUploader struct {
	rs              *requestsender.RequestSender
	ub              urlBuilder
	repeatIntervals []int
}

// Execute выполняет загрузку метрик
func (c *jsonMetricsUploader) Execute(result dto.MetricsDto) error {
	var list []dto.JSONMetrics
	for _, metric := range result.All() {
		list = append(list, metric)
	}
	if len(list) == 0 {
		return nil
	}
	err := c.sendMany(list)
	if err != nil {
		logger.Log.Error("error send response", zap.Error(err))
	}
	return nil
}

func (c *jsonMetricsUploader) Process(results <-chan dto.MetricsDto) error {
	for result := range results {
		err := c.Execute(result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *jsonMetricsUploader) send(m dto.JSONMetrics) error {
	url := c.ub.Build(types.URLParams{})
	encoded, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return c.rs.Post(url, contenttypes.ApplicationJSON, encoded)
}

func (c *jsonMetricsUploader) sendMany(l []dto.JSONMetrics) error {
	url := c.ub.Build(types.URLParams{"endpoint": "updates"})
	encoded, err := json.Marshal(l)
	if err != nil {
		return err
	}
	err = c.rs.Post(url, contenttypes.ApplicationJSON, encoded)
	if err == nil {
		return nil
	}
	for _, repeatInterval := range c.repeatIntervals {
		time.Sleep(time.Duration(repeatInterval) * time.Second)
		err = c.rs.Post(url, contenttypes.ApplicationJSON, encoded)
		if err == nil {
			return nil
		}
	}
	return err
}

// NewJSONUploader возвращает сервис jsonMetricsUploader для загрузки метрик в формате JSON
func NewJSONUploader(rs *requestsender.RequestSender, ub urlBuilder, repeatIntervals []int) *jsonMetricsUploader {
	return &jsonMetricsUploader{rs, ub, repeatIntervals}
}
