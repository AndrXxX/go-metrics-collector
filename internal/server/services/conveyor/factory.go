package conveyor

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/stack"
)

type conveyorFactory struct {
	l logger
}

func (f *conveyorFactory) From(list []interfaces.Handler) *handlersConveyor {
	return &handlersConveyor{
		logger: f.l,
		stack:  *stack.NewFromSlice(list),
	}
}

func Factory(l logger) *conveyorFactory {
	return &conveyorFactory{l}
}
