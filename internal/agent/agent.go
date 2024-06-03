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
	e := getExecutors(config)
	sleepInterval := getSleepInterval(config)
	for {
		for _, executor := range e {
			err := executor.Execute(*m)
			if err != nil {
				return err
			}
		}
		time.Sleep(sleepInterval)
	}
}

func getExecutors(config *config.Config) []executors.Executors {
	list := make([]executors.Executors, 0)
	list = append(list, executors.NewCollector(time.Duration(config.Intervals.PollInterval), config))
	list = append(list, executors.NewUploader(time.Duration(config.Intervals.ReportInterval), config))
	return list
}

func getSleepInterval(config *config.Config) time.Duration {
	return time.Duration(math.Min(float64(config.Intervals.PollInterval), float64(config.Intervals.ReportInterval))) * time.Second
}
