package metricsformatter

import (
	"errors"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsValueFormatter struct {
}

func (s MetricsValueFormatter) Format(m *models.Metrics) (string, error) {
	switch m.MType {
	case metrics.Counter:
		if m.Delta == nil {
			return "", nil
		}
		return fmt.Sprintf("%d", *m.Delta), nil
	case metrics.Gauge:
		if m.Value == nil {
			return "", nil
		}
		return fmt.Sprintf("%v", *m.Value), nil
	default:
		return "", errors.New("unknown metric type")
	}
}
