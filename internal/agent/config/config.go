package config

type Config struct {
	Common    CommonConfig
	Intervals Intervals
	Metrics   MetricsList
}

type MetricsList []string

type Intervals struct {
	PollInterval   int64 `valid:"range(1|999)"`
	ReportInterval int64 `valid:"range(1|999)"`
	SleepInterval  int64 `valid:"range(1|20)"`
}

type CommonConfig struct {
	Host     string `valid:"minstringlength(3)"`
	LogLevel string `valid:"in(debug|info|warn|error|fatal)"`
}
