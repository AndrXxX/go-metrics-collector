package runtimemetricscollector

import (
	"sync"
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
	tests := []struct {
		name    string
		ml      *config.MetricsList
		result  dto.MetricsDto
		include []string
		exclude []string
	}{
		{
			name:    "Test Unknown field",
			ml:      &config.MetricsList{"UnknownMetric"},
			result:  *dto.NewMetricsDto(),
			include: []string{metrics.PollCount},
			exclude: []string{"UnknownMetric"},
		},
		{
			name:    "Test Counter",
			ml:      &config.MetricsList{},
			result:  *dto.NewMetricsDto(),
			include: []string{metrics.PollCount},
		},
		{
			name: "Test Counter when exist",
			ml:   &config.MetricsList{},
			result: func() dto.MetricsDto {
				v := *dto.NewMetricsDto()
				var current int64 = 1
				v.Set(dto.JSONMetrics{ID: metrics.PollCount, Delta: &current})
				return v
			}(),
			include: []string{metrics.PollCount},
		},
		{
			name:    "Test Alloc",
			ml:      &config.MetricsList{metrics.Alloc},
			result:  *dto.NewMetricsDto(),
			include: []string{metrics.Alloc, metrics.RandomValue},
		},
		{
			name:    "Test BuckHashSys, HeapObjects",
			ml:      &config.MetricsList{metrics.BuckHashSys, metrics.HeapObjects},
			result:  *dto.NewMetricsDto(),
			include: []string{metrics.BuckHashSys, metrics.HeapObjects},
		},
		{
			name:    "Test empty",
			ml:      &config.MetricsList{},
			result:  *dto.NewMetricsDto(),
			include: []string{metrics.RandomValue},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collector{ml: tt.ml}
			c.execute(tt.result)
			for _, v := range tt.include {
				_, ok := tt.result.Get(v)
				assert.True(t, ok)
			}
			if tt.exclude != nil {
				for _, v := range tt.exclude {
					_, ok := tt.result.Get(v)
					assert.False(t, ok)
				}
			}
		})
	}
}

func Test_collector_Collect(t *testing.T) {
	tests := []struct {
		name    string
		ml      *config.MetricsList
		include []string
	}{
		{
			name:    "Test Counter",
			ml:      &config.MetricsList{},
			include: []string{metrics.PollCount},
		},
		{
			name:    "Test Alloc",
			ml:      &config.MetricsList{metrics.Alloc},
			include: []string{metrics.Alloc, metrics.RandomValue},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan dto.MetricsDto)
			c := New(tt.ml)
			wg.Add(1)
			go func() {
				err := c.Collect(ch)
				wg.Done()
				assert.NoError(t, err)
			}()
			wg.Add(1)
			go func() {
				for result := range ch {
					for _, v := range tt.include {
						_, ok := result.Get(v)
						assert.True(t, ok)
					}
				}
				wg.Done()
			}()
			wg.Wait()
		})
	}
}
