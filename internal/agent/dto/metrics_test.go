package dto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *MetricsDto
	}{
		{
			name: "Test New MemStorage",
			want: &MetricsDto{Gauge: map[string]float64{}, Counter: map[string]int64{}, list: map[string]JSONMetrics{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetricsDto()
			assert.Equal(t, tt.want, m)
		})
	}
}
