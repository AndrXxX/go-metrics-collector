package metricschecker

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type metricsChecker struct {
	validTypes map[string]struct{}
}

func (c *metricsChecker) IsValid(m *models.Metrics) bool {
	_, ok := c.validTypes[m.MType]
	return ok
}

func New() *metricsChecker {
	return &metricsChecker{
		validTypes: map[string]struct{}{
			metrics.Counter: {},
			metrics.Gauge:   {},
		},
	}
}
