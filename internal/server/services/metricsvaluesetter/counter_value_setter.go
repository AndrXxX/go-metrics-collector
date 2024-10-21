package metricsvaluesetter

import (
	"fmt"
	"strconv"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type counterValueSetter struct {
}

// Set устанавливает значение метрики типа counter
func (i *counterValueSetter) Set(m *models.Metrics, value string) error {
	if value == "" {
		return fmt.Errorf("empty value for couner metric")
	}
	converted, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	m.Delta = &converted
	return nil
}
