package config

type Config struct {
	Common    CommonConfig
	Intervals Intervals
	Metrics   MetricsList
}

type MetricsList []string

type Intervals struct {
	PollInterval   int64
	ReportInterval int64
}

type CommonConfig struct {
	Host string
}
