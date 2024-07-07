package metricsvaluesetter

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"strconv"
)

type gaugeValueSetter struct {
}

func (i *gaugeValueSetter) Set(m *models.Metrics, value string) error {
	if value == "" {
		return fmt.Errorf("empty value for gauge metric")
	}
	converted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.Value = &converted
	return nil
}
