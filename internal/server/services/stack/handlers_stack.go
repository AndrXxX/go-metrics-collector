package stack

import "net/http"

type handler interface {
	Execute(http.ResponseWriter, *http.Request) (ok bool)
}

func ForHandlers(list []handler) *stack[handler] {
	return NewFromSlice(list)
}
