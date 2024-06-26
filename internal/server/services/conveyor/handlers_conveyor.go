package conveyor

import "net/http"

type stackInterface[T any] interface {
	Push(value T)
	Pop() (T, bool)
}

type handler interface {
	Execute(http.ResponseWriter, *http.Request) (ok bool)
}

type handlersConveyor struct {
	stack stackInterface[handler]
}

func (c *handlersConveyor) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for {
			handler, ok := c.stack.Pop()
			if !ok || !handler.Execute(w, r) {
				return
			}
		}
	}
}

func New(stack stackInterface[handler]) *handlersConveyor {
	return &handlersConveyor{stack: stack}
}
