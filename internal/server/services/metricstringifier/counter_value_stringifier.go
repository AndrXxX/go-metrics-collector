package metricstringifier

import (
	"errors"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsValueStringifier struct {
}

func (s MetricsValueStringifier) String(m *models.Metrics) (string, error) {
	switch m.MType {
	case metrics.Counter:
		return fmt.Sprintf("%d", m.Delta), nil
	case metrics.Gauge:
		return fmt.Sprintf("%d", m.Value), nil
	default:
		return "", errors.New("unknown metric type")
	}
}
