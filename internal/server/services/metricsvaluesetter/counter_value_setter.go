package metricsvaluesetter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"strconv"
)

type counterValueSetter struct {
}

func (i *counterValueSetter) Set(m *models.Metrics, value string) error {
	converted, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	m.Delta = &converted
	return nil
}
