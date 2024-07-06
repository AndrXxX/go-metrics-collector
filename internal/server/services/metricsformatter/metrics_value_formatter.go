package metricsformatter

import (
	"errors"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"strconv"
)

type MetricsValueFormatter struct {
}

func (s MetricsValueFormatter) Format(m *models.Metrics) (string, error) {
	switch m.MType {
	case metrics.Counter:
		if m.Delta == nil {
			return "", nil
		}
		return strconv.FormatInt(*m.Delta, 10), nil
	case metrics.Gauge:
		if m.Value == nil {
			return "", nil
		}
		return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
	default:
		return "", errors.New("unknown metric type")
	}
}
