package executors

import "github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"

type Executors interface {
	Execute(metrics.Metrics) error
}
