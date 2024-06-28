package executors

import "github.com/AndrXxX/go-metrics-collector/internal/agent/dto"

type Executor interface {
	Execute(dto.MetricsDto) error
}
