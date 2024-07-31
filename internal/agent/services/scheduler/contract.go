package scheduler

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"time"
)

type collector interface {
	Collect(chan<- dto.MetricsDto) error
}

type processor interface {
	Process(<-chan dto.MetricsDto) error
}

type collectorItem struct {
	c            collector
	interval     time.Duration
	lastExecuted time.Time
}

type processorItem struct {
	p            processor
	interval     time.Duration
	lastExecuted time.Time
}
