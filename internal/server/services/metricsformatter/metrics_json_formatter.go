package metricsformatter

import (
	"encoding/json"
	"fmt"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsJSONFormatter struct {
}

func (s MetricsJSONFormatter) Format(m *models.Metrics) (string, error) {
	res, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("error on marshalling metrics: %w", err)
	}
	return string(res), err
}
