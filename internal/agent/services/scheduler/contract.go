package scheduler

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"time"
)

type executor interface {
	Execute(metrics.Metrics) error
}

type item struct {
	e            executor
	interval     time.Duration
	lastExecuted time.Time
}
