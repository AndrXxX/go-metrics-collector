package metricsuploader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
)

type jsonMetricsUploader struct {
	rs              requestSender
	ub              urlBuilder
	repeatIntervals []int
}

func (c *jsonMetricsUploader) execute(result dto.MetricsDto) error {
	var list []dto.JSONMetrics
	for _, metric := range result.All() {
		list = append(list, metric)
	}
	if len(list) == 0 {
		return nil
	}
	err := c.sendMany(list)
	if err != nil {
		return fmt.Errorf("error send response: %w", err)
	}
	return nil
}

// Process выполняет загрузку метрик
func (c *jsonMetricsUploader) Process(results <-chan dto.MetricsDto) error {
	for result := range results {
		err := c.execute(result)
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
func NewJSONUploader(rs requestSender, ub urlBuilder, repeatIntervals []int) *jsonMetricsUploader {
	return &jsonMetricsUploader{rs, ub, repeatIntervals}
}
