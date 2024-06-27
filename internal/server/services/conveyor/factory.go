package conveyor

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/stack"
)

type conveyorFactory struct {
	logger loggerFunc
}

func (f *conveyorFactory) From(list []interfaces.Handler) *handlersConveyor {
	return &handlersConveyor{
		logger: f.logger,
		stack:  *stack.NewFromSlice(list),
	}
}

func Factory(logger loggerFunc) *conveyorFactory {
	return &conveyorFactory{
		logger: logger,
	}
}
