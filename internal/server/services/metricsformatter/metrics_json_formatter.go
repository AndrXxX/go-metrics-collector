package metricsformatter

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

// MetricsJSONFormatter сервис для форматирования метрики в формат JSON
type MetricsJSONFormatter struct {
}

// Format возвращает строку в формате JSON
func (s MetricsJSONFormatter) Format(m *models.Metrics) (string, error) {
	res, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("error on marshalling metrics: %w", err)
	}
	return string(res), err
}
