package metricstringifier

import (
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsJSONStringifier struct {
}

func (s MetricsJSONStringifier) String(m *models.Metrics) (string, error) {
	res, err := json.Marshal(m)
	return string(res), err
}
