package metricschecker

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *metricsChecker
	}{
		{
			name: "OK Test",
			want: &metricsChecker{map[string]struct{}{metrics.Counter: {}, metrics.Gauge: {}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New())
		})
	}
}

func TestMetricsCheckerIsValid(t *testing.T) {
	tests := []struct {
		name string
		mc   *metricsChecker
		m    *models.Metrics
		want bool
	}{
		{
			name: "Valid Counter Test",
			mc:   &metricsChecker{map[string]struct{}{metrics.Counter: {}, metrics.Gauge: {}}},
			m:    &models.Metrics{MType: metrics.Counter},
			want: true,
		},
		{
			name: "Valid Gauge Test",
			mc:   &metricsChecker{map[string]struct{}{metrics.Counter: {}, metrics.Gauge: {}}},
			m:    &models.Metrics{MType: metrics.Gauge},
			want: true,
		},
		{
			name: "Not valid Test",
			mc:   &metricsChecker{map[string]struct{}{metrics.Counter: {}, metrics.Gauge: {}}},
			m:    &models.Metrics{MType: "not valid"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mc.IsValid(tt.m))
		})
	}
}
