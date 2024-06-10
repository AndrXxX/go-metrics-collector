package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	me "github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollector(t *testing.T) {
	tests := []struct {
		name string
		ml   *config.MetricsList
		want Executors
	}{
		{
			name: "Test New metricsCollector #1 (Alloc)",
			ml:   &config.MetricsList{me.Alloc},
			want: &metricsCollector{ml: &config.MetricsList{me.Alloc}},
		},
		{
			name: "Test New metricsCollector #1 (Alloc)",
			ml:   &config.MetricsList{me.BuckHashSys, me.HeapObjects},
			want: &metricsCollector{ml: &config.MetricsList{me.BuckHashSys, me.HeapObjects}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector(tt.ml)
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
		result metrics.Metrics
		valuesInResult
	}{
		{
			name:   "Test Counter",
			ml:     &config.MetricsList{me.Alloc},
			result: *metrics.NewMetrics(),
			valuesInResult: valuesInResult{
				counter: []string{me.PollCount},
			},
		},
		{
			name:   "Test Alloc",
			ml:     &config.MetricsList{me.Alloc},
			result: *metrics.NewMetrics(),
			valuesInResult: valuesInResult{
				gauge: []string{me.Alloc, me.RandomValue},
			},
		},
		{
			name:   "Test BuckHashSys, HeapObjects",
			ml:     &config.MetricsList{me.BuckHashSys, me.HeapObjects},
			result: *metrics.NewMetrics(),
			valuesInResult: valuesInResult{
				gauge: []string{me.BuckHashSys, me.HeapObjects},
			},
		},
		{
			name:   "Test empty",
			ml:     &config.MetricsList{},
			result: *metrics.NewMetrics(),
			valuesInResult: valuesInResult{
				gauge: []string{me.RandomValue},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &metricsCollector{ml: tt.ml}
			err := c.Execute(tt.result)
			assert.NoError(t, err)
			for _, v := range tt.valuesInResult.gauge {
				assert.NotNil(t, tt.result.Gauge[v])
			}
			for _, v := range tt.valuesInResult.counter {
				assert.NotNil(t, tt.result.Counter[v])
			}
		})
	}
}
