package scheduler

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
)

type Collector interface {
	Collect(chan<- dto.MetricsDto) error
}

type Processor interface {
	Process(<-chan dto.MetricsDto) error
}
