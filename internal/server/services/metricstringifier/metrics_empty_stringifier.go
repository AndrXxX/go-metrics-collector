package metricstringifier

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsEmptyStringifier struct {
}

func (s MetricsEmptyStringifier) String(_ *models.Metrics) (string, error) {
	return "", nil
}
