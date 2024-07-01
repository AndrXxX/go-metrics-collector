package metricsvaluesetter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"strconv"
)

type GaugeValueSetter struct {
}

func (i *GaugeValueSetter) Set(m *models.Metrics, value string) error {
	converted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.Value = &converted
	return nil
}
