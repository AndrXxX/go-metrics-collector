package metricschecker

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"slices"
)

type metricsChecker struct {
	validTypes []string
}

func (c *metricsChecker) IsValid(m *models.Metrics) bool {
	return slices.Contains(c.validTypes, m.MType)
}

func New() *metricsChecker {
	return &metricsChecker{
		validTypes: []string{metrics.Counter, metrics.Gauge},
	}
}
