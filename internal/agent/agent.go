package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/executors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"math"
	"time"
)

func Run(config *config.Config) error {
	m := metrics.NewMetrics()
	e := getExecutors(config, m)
	sleepInterval := getSleepInterval(config)
	for {
		for _, executor := range e {
			err := executor.Execute()
			if err != nil {
				return err
			}
		}
		time.Sleep(sleepInterval)
	}
}

func getExecutors(config *config.Config, result *metrics.Metrics) []executors.Executors {
	list := make([]executors.Executors, 0)
	list = append(list, executors.NewCollector(time.Duration(config.Intervals.PollInterval), config, result))
	list = append(list, executors.NewUploader(time.Duration(config.Intervals.ReportInterval), config, result))
	return list
}

func getSleepInterval(config *config.Config) time.Duration {
	return time.Duration(math.Min(float64(config.Intervals.PollInterval), float64(config.Intervals.ReportInterval))) * time.Second
}
