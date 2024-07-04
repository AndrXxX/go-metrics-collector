package metricsformatter

import (
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsJSONFormatter struct {
}

func (s MetricsJSONFormatter) Format(m *models.Metrics) (string, error) {
	res, err := json.Marshal(m)
	return string(res), err
}
