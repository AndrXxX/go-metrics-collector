package config

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars/defaults"
)

func NewConfig() *Config {
	return &Config{
		Common: CommonConfig{
			Host:      defaults.Host,
			LogLevel:  defaults.LogLevel,
			Key:       defaults.Key,
			RateLimit: defaults.RateLimit,
		},
		Intervals: Intervals{
			PollInterval:    defaults.PollInterval,
			ReportInterval:  defaults.ReportInterval,
			SleepInterval:   defaults.SleepInterval,
			RepeatIntervals: defaults.RepeatIntervals,
		},
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
