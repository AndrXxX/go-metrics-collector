package executors

import "github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"

type Executor interface {
	Execute(metrics.Metrics) error
}
