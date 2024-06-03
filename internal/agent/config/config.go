package config

import me "github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"

type Config struct {
	Common    CommonConfig
	Intervals Intervals
	Metrics   []string
}

type Intervals struct {
	PollInterval   int64
	ReportInterval int64
}

type CommonConfig struct {
	Host string
}

func NewConfig() *Config {
	return &Config{
		Common: CommonConfig{
			Host: "http://localhost:8080",
		},
		Intervals: Intervals{PollInterval: 2, ReportInterval: 10},
		Metrics: []string{
			me.Alloc,
			me.BuckHashSys,
			me.Frees,
			me.GCCPUFraction,
			me.GCSys,
			me.HeapAlloc,
			me.HeapIdle,
			me.HeapInuse,
			me.HeapObjects,
			me.HeapReleased,
			me.HeapSys,
			me.LastGC,
			me.Lookups,
			me.MCacheInuse,
			me.MCacheSys,
			me.MSpanInuse,
			me.MSpanSys,
			me.Mallocs,
			me.NextGC,
			me.NumForcedGC,
			me.NumGC,
			me.OtherSys,
			me.PauseTotalNs,
			me.StackInuse,
			me.StackSys,
			me.Sys,
			me.TotalAlloc,
		},
	}
}
