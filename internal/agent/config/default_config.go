package config

import "github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"

func NewConfig() *Config {
	return &Config{
		Common: CommonConfig{
			Host: "localhost:8080",
		},
		Intervals: Intervals{PollInterval: 2, ReportInterval: 10, SleepInterval: 1},
		Metrics: MetricsList{
			metrics.Alloc,
			metrics.BuckHashSys,
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
	}
}
