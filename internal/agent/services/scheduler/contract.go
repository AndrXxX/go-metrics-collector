package scheduler

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
)

type collector interface {
	Collect(chan<- dto.MetricsDto) error
}

type processor interface {
	Process(<-chan dto.MetricsDto) error
}
