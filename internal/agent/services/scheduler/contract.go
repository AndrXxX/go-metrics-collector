package scheduler

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"time"
)

type executor interface {
	Execute(dto.MetricsDto) error
}

type item struct {
	e            executor
	interval     time.Duration
	lastExecuted time.Time
}
