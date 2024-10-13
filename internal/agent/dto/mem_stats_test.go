package dto

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

func TestNewMemStatsDto(t *testing.T) {
	memStats := &runtime.MemStats{}
	tests := []struct {
		name string
		want memStatsDto
	}{
		{
			name: "OK",
			want: memStatsDto{gettersMap: buildMap(memStats)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ObjectsAreEqual(tt.want, NewMemStatsDto(memStats))
		})
	}
}

func Test_buildMap(t *testing.T) {
	tests := []struct {
		name    string
		metrics []string
	}{
		{
			name:    "Exist Alloc, BuckHashSys",
			metrics: []string{metrics.Alloc, metrics.BuckHashSys},
		},
		{
			name: "Exist all others",
			metrics: []string{
				metrics.Frees,
				metrics.GCCPUFraction,
				metrics.GCSys,
				metrics.HeapAlloc,
				metrics.HeapIdle,
				metrics.HeapInuse,
				metrics.HeapObjects,
				metrics.HeapReleased,
				metrics.HeapSys,
				metrics.LastGC,
				metrics.Lookups,
				metrics.MCacheInuse,
				metrics.MCacheSys,
				metrics.MSpanInuse,
				metrics.MSpanSys,
				metrics.Mallocs,
				metrics.NextGC,
				metrics.NumForcedGC,
				metrics.NumGC,
				metrics.OtherSys,
				metrics.PauseTotalNs,
				metrics.StackInuse,
				metrics.StackSys,
				metrics.Sys,
				metrics.TotalAlloc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsMap := buildMap(&runtime.MemStats{})
			for _, m := range tt.metrics {
				if assert.Contains(t, metricsMap, m) {
					assert.IsType(t, float64(0), metricsMap[m]())
				}
			}
		})
	}
}

func Test_memStatsDto_FetchGetter(t *testing.T) {
	ms := NewMemStatsDto(&runtime.MemStats{})
	tests := []struct {
		name      string
		metric    string
		want      getter
		wantExist bool
	}{
		{
			name:      "Test with exist metric Alloc",
			metric:    metrics.Alloc,
			want:      ms.gettersMap[metrics.Alloc],
			wantExist: true,
		},
		{
			name:      "Test with not exist metric unknown",
			metric:    "unknown",
			want:      ms.gettersMap["unknown"],
			wantExist: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			getterFunc, exist := ms.FetchGetter(tt.metric)
			assert.Equal(t, tt.wantExist, exist)
			assert.Equal(t, tt.want != nil, getterFunc != nil)
		})
	}
}
