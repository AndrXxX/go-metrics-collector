package dto

import (
	"runtime"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

type getter func() float64

type memStatsDto struct {
	gettersMap map[string]getter
}

func (ms *memStatsDto) FetchGetter(name string) (getter, bool) {
	getter, ok := ms.gettersMap[name]
	return getter, ok
}

func NewMemStatsDto(ms *runtime.MemStats) memStatsDto {
	return memStatsDto{gettersMap: buildMap(ms)}
}

func buildMap(s *runtime.MemStats) map[string]getter {
	m := make(map[string]getter)
	m[metrics.Alloc] = func() float64 { return float64(s.Alloc) }
	m[metrics.BuckHashSys] = func() float64 { return float64(s.BuckHashSys) }
	m[metrics.Frees] = func() float64 { return float64(s.Frees) }
	m[metrics.GCCPUFraction] = func() float64 { return s.GCCPUFraction }
	m[metrics.GCSys] = func() float64 { return float64(s.GCSys) }
	m[metrics.HeapAlloc] = func() float64 { return float64(s.HeapAlloc) }
	m[metrics.HeapIdle] = func() float64 { return float64(s.HeapIdle) }
	m[metrics.HeapInuse] = func() float64 { return float64(s.HeapInuse) }
	m[metrics.HeapObjects] = func() float64 { return float64(s.HeapObjects) }
	m[metrics.HeapReleased] = func() float64 { return float64(s.HeapReleased) }
	m[metrics.HeapSys] = func() float64 { return float64(s.HeapSys) }
	m[metrics.LastGC] = func() float64 { return float64(s.LastGC) }
	m[metrics.Lookups] = func() float64 { return float64(s.Lookups) }
	m[metrics.MCacheInuse] = func() float64 { return float64(s.MCacheInuse) }
	m[metrics.MCacheSys] = func() float64 { return float64(s.MCacheSys) }
	m[metrics.MSpanInuse] = func() float64 { return float64(s.MSpanInuse) }
	m[metrics.MSpanSys] = func() float64 { return float64(s.MSpanSys) }
	m[metrics.Mallocs] = func() float64 { return float64(s.Mallocs) }
	m[metrics.NextGC] = func() float64 { return float64(s.NextGC) }
	m[metrics.NumForcedGC] = func() float64 { return float64(s.NumForcedGC) }
	m[metrics.NumGC] = func() float64 { return float64(s.NumGC) }
	m[metrics.OtherSys] = func() float64 { return float64(s.OtherSys) }
	m[metrics.PauseTotalNs] = func() float64 { return float64(s.PauseTotalNs) }
	m[metrics.StackInuse] = func() float64 { return float64(s.StackInuse) }
	m[metrics.StackSys] = func() float64 { return float64(s.StackSys) }
	m[metrics.Sys] = func() float64 { return float64(s.Sys) }
	m[metrics.TotalAlloc] = func() float64 { return float64(s.TotalAlloc) }
	return m
}
