package metricsformatter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

// MetricsEmptyFormatter сервис для форматирования метрики в пустое значение
type MetricsEmptyFormatter struct {
}

// Format возвращает пустое значение
func (s MetricsEmptyFormatter) Format(_ *models.Metrics) (string, error) {
	return "", nil
}
