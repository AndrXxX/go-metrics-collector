package utils

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"runtime"
)

type extendedMemStats struct {
	Stats runtime.MemStats
}

func NewExtendedMemStats() extendedMemStats {
	return extendedMemStats{Stats: runtime.MemStats{}}
}

func (ms *extendedMemStats) GetValue(name string) (float64, error) {
	switch name {
	case metrics.Alloc:
		return float64(ms.Stats.Alloc), nil
	case metrics.BuckHashSys:
		return float64(ms.Stats.BuckHashSys), nil
	case metrics.Frees:
		return float64(ms.Stats.Frees), nil
	case metrics.GCCPUFraction:
		return ms.Stats.GCCPUFraction, nil
	case metrics.GCSys:
		return float64(ms.Stats.GCSys), nil
	case metrics.HeapAlloc:
		return float64(ms.Stats.HeapAlloc), nil
	case metrics.HeapIdle:
		return float64(ms.Stats.HeapIdle), nil
	case metrics.HeapInuse:
		return float64(ms.Stats.HeapInuse), nil
	case metrics.HeapObjects:
		return float64(ms.Stats.HeapObjects), nil
	case metrics.HeapReleased:
		return float64(ms.Stats.HeapReleased), nil
	case metrics.HeapSys:
		return float64(ms.Stats.HeapSys), nil
	case metrics.LastGC:
		return float64(ms.Stats.LastGC), nil
	case metrics.Lookups:
		return float64(ms.Stats.Lookups), nil
	case metrics.MCacheInuse:
		return float64(ms.Stats.MCacheInuse), nil
	case metrics.MCacheSys:
		return float64(ms.Stats.MCacheSys), nil
	case metrics.MSpanInuse:
		return float64(ms.Stats.MSpanInuse), nil
	case metrics.MSpanSys:
		return float64(ms.Stats.MSpanSys), nil
	case metrics.Mallocs:
		return float64(ms.Stats.Mallocs), nil
	case metrics.NextGC:
		return float64(ms.Stats.NextGC), nil
	case metrics.NumForcedGC:
		return float64(ms.Stats.NumForcedGC), nil
	case metrics.NumGC:
		return float64(ms.Stats.NumGC), nil
	case metrics.OtherSys:
		return float64(ms.Stats.OtherSys), nil
	case metrics.PauseTotalNs:
		return float64(ms.Stats.PauseTotalNs), nil
	case metrics.StackInuse:
		return float64(ms.Stats.StackInuse), nil
	case metrics.StackSys:
		return float64(ms.Stats.StackSys), nil
	case metrics.Sys:
		return float64(ms.Stats.Sys), nil
	case metrics.TotalAlloc:
		return float64(ms.Stats.TotalAlloc), nil
	default:
		return 0, fmt.Errorf("no such metric")
	}
}
