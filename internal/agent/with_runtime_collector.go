package agent

import "github.com/AndrXxX/go-metrics-collector/internal/agent/services/runtimemetricscollector"

func WithRuntimeCollector() Option {
	return func(a agent) {
		a.collectors.Add(runtimemetricscollector.New(&a.c.Metrics))
	}
}
