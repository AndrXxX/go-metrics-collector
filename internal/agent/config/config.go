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
	SleepInterval  int64
	// TODO:  Perederey 2 weeks ago Лучше использовать struct tags для валидации значений конфигурации.
}

type CommonConfig struct {
	Host     string
	LogLevel string
}
