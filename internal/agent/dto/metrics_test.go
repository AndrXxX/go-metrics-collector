package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *MetricsDto
	}{
		{
			name: "Test NewMetricsDto",
			want: &MetricsDto{list: map[string]JSONMetrics{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetricsDto()
			assert.Equal(t, tt.want, m)
		})
	}
}

func TestMetricsDto_All(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]JSONMetrics
		want   map[string]JSONMetrics
	}{
		{
			name: "Test OK with two metrics",
			fields: map[string]JSONMetrics{
				metrics.Alloc:     {ID: metrics.Alloc},
				metrics.HeapAlloc: {ID: metrics.HeapAlloc},
			},
			want: map[string]JSONMetrics{
				metrics.Alloc:     {ID: metrics.Alloc},
				metrics.HeapAlloc: {ID: metrics.HeapAlloc},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &MetricsDto{
				list: tt.fields,
			}
			assert.Equal(t, tt.want, dto.All())
		})
	}
}

func TestMetricsDto_Get(t *testing.T) {
	tests := []struct {
		name      string
		fields    map[string]JSONMetrics
		metric    string
		want      JSONMetrics
		wantExist bool
	}{
		{
			name: "Test exist with Alloc metric",
			fields: map[string]JSONMetrics{
				metrics.Alloc: {ID: metrics.Alloc},
			},
			metric:    metrics.Alloc,
			want:      JSONMetrics{ID: metrics.Alloc},
			wantExist: true,
		},
		{
			name: "Test not exist with unknown metric",
			fields: map[string]JSONMetrics{
				metrics.Alloc: {ID: metrics.Alloc}},
			metric:    "unknown",
			want:      JSONMetrics{},
			wantExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &MetricsDto{
				list: tt.fields,
			}

			metric, exist := dto.Get(tt.metric)
			assert.Equal(t, tt.want, metric)
			assert.Equal(t, tt.wantExist, exist)
		})
	}
}

func TestMetricsDto_Set(t *testing.T) {
	getPointer := func(val float64) *float64 { return &val }
	tests := []struct {
		name   string
		fields map[string]JSONMetrics
		metric JSONMetrics
	}{
		{
			name:   "Set Alloc metric when not exist",
			fields: map[string]JSONMetrics{},
			metric: JSONMetrics{ID: metrics.Alloc},
		},
		{
			name:   "Set Alloc metric when exist",
			fields: map[string]JSONMetrics{metrics.Alloc: {ID: metrics.Alloc, Value: getPointer(9)}},
			metric: JSONMetrics{ID: metrics.Alloc, Value: getPointer(10.1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &MetricsDto{
				list: tt.fields,
			}
			dto.Set(tt.metric)
			assert.Equal(t, tt.metric, dto.list[metrics.Alloc])
		})
	}
}
