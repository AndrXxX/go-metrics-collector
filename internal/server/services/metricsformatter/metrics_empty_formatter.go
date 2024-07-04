package metricsformatter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsEmptyFormatter struct {
}

func (s MetricsEmptyFormatter) Format(_ *models.Metrics) (string, error) {
	return "", nil
}
