package conveyor

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/stack"
	"net/http"
)

type handlersConveyor struct {
	stack  stack.Stack[interfaces.Handler]
	logger logger
}

func (c *handlersConveyor) Handler() http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		newStack := c.stack.Copy()
		var next http.HandlerFunc
		next = func(w http.ResponseWriter, r *http.Request) {
			h, ok := newStack.Shift()
			if !ok || h == nil {
				return
			}
			h.Handle(w, r, next)
		}
		next(w, r)
	}
	if c.logger != nil {
		return c.logger.Handle(handler)
	}
	return handler
}

func (c *handlersConveyor) Add(handler interfaces.Handler) *handlersConveyor {
	c.stack.Push(handler)
	return c
}
