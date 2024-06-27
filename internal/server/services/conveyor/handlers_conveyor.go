package conveyor

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/stack"
	"net/http"
)

type loggerFunc func(h http.HandlerFunc) http.HandlerFunc

type handlersConveyor struct {
	stack  stack.Stack[interfaces.Handler]
	logger loggerFunc
}

func (c *handlersConveyor) Handler() http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		for {
			handler, ok := c.stack.Pop()
			if !ok || !handler.Handle(w, r) {
				return
			}
		}
	}
	if c.logger != nil {
		return c.logger(handler)
	}
	return handler
}

func (c *handlersConveyor) Add(handler interfaces.Handler) *handlersConveyor {
	c.stack.Push(handler)
	return c
}

func (c *handlersConveyor) From(list []interfaces.Handler) *handlersConveyor {
	c.stack = *stack.NewFromSlice(list)
	return c
}

func New(logger loggerFunc) *handlersConveyor {
	return &handlersConveyor{
		stack:  *stack.NewFromSlice([]interfaces.Handler{}),
		logger: logger,
	}
}
