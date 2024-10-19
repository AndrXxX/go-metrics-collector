package runtimemetricscollector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

func TestNewCollector(t *testing.T) {
	tests := []struct {
		name string
		ml   *config.MetricsList
		want *collector
	}{
		{
			name: "Test New collector #1 (Alloc)",
			ml:   &config.MetricsList{metrics.Alloc},
			want: &collector{ml: &config.MetricsList{metrics.Alloc}},
		},
		{
			name: "Test New collector #1 (Alloc)",
			ml:   &config.MetricsList{metrics.BuckHashSys, metrics.HeapObjects},
			want: &collector{ml: &config.MetricsList{metrics.BuckHashSys, metrics.HeapObjects}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.ml)
			assert.Equal(t, tt.want, c)
		})
	}
}

func Test_metricsCollector_Execute(t *testing.T) {
	type valuesInResult struct {
		gauge   []string
		counter []string
	}

	tests := []struct {
		name   string
		ml     *config.MetricsList
		result dto.MetricsDto
		valuesInResult
		excludeValuesInResult valuesInResult
	}{
		{
			name:   "Test Unknown field",
			ml:     &config.MetricsList{"UnknownMetric"},
			result: *dto.NewMetricsDto(),
			valuesInResult: valuesInResult{
				counter: []string{metrics.PollCount},
			},
			excludeValuesInResult: valuesInResult{
				counter: []string{"UnknownMetric"},
				gauge:   []string{"UnknownMetric"},
			},
		},
		{
			name:   "Test Counter",
			ml:     &config.MetricsList{},
			result: *dto.NewMetricsDto(),
			valuesInResult: valuesInResult{
				counter: []string{metrics.PollCount},
			},
		},
		{
			name:   "Test Alloc",
			ml:     &config.MetricsList{metrics.Alloc},
			result: *dto.NewMetricsDto(),
			valuesInResult: valuesInResult{
				gauge: []string{metrics.Alloc, metrics.RandomValue},
			},
		},
		{
			name:   "Test BuckHashSys, HeapObjects",
			ml:     &config.MetricsList{metrics.BuckHashSys, metrics.HeapObjects},
			result: *dto.NewMetricsDto(),
			valuesInResult: valuesInResult{
				gauge: []string{metrics.BuckHashSys, metrics.HeapObjects},
			},
		},
		{
			name:   "Test empty",
			ml:     &config.MetricsList{},
			result: *dto.NewMetricsDto(),
			valuesInResult: valuesInResult{
				gauge: []string{metrics.RandomValue},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collector{ml: tt.ml}
			err := c.execute(tt.result)
			assert.NoError(t, err)
			for _, v := range tt.valuesInResult.gauge {
				_, ok := tt.result.Get(v)
				assert.True(t, ok)
			}
			for _, v := range tt.valuesInResult.counter {
				_, ok := tt.result.Get(v)
				assert.True(t, ok)
			}
			if tt.excludeValuesInResult.gauge != nil {
				for _, v := range tt.excludeValuesInResult.gauge {
					_, ok := tt.result.Get(v)
					assert.False(t, ok)
				}
				for _, v := range tt.excludeValuesInResult.counter {
					_, ok := tt.result.Get(v)
					assert.False(t, ok)
				}
			}
		})
	}
}
