package interfaces

import "net/http"

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type Middleware interface {
	Handler(next http.Handler) http.Handler
}
