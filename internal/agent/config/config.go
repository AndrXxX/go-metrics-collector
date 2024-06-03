package config

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
			"Alloc",
			"BuckHashSys",
			"Frees",
			"GCCPUFraction",
			"GCSys",
			"HeapAlloc",
			"HeapIdle",
			"HeapInuse",
			"HeapObjects",
			"HeapReleased",
			"HeapSys",
			"LastGC",
			"Lookups",
			"MCacheInuse",
			"MCacheSys",
			"MSpanInuse",
			"MSpanSys",
			"Mallocs",
			"NextGC",
			"NumForcedGC",
			"NumGC",
			"OtherSys",
			"PauseTotalNs",
			"StackInuse",
			"StackSys",
			"Sys",
			"TotalAlloc",
		},
	}
}
