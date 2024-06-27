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
		for _, handler := range c.stack.All() {
			if !handler.Handle(w, r) {
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
