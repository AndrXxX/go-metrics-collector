package config

// Config конфигурация агента
type Config struct {
	Common    CommonConfig
	Intervals Intervals
	Metrics   MetricsList
}

// MetricsList настройка списка метрик
type MetricsList []string

// Intervals настройки интервалов
type Intervals struct {
	PollInterval    int64 `valid:"required,range(1|999)"`
	ReportInterval  int64 `valid:"required,range(1|999)"`
	SleepInterval   int64 `valid:"required,range(1|20)"`
	RepeatIntervals []int
}

// CommonConfig общая конфигурация агента
type CommonConfig struct {
	Host      string `valid:"minstringlength(3)"`
	LogLevel  string `valid:"in(debug|info|warn|error|fatal)"`
	Key       string
	RateLimit int64
}
