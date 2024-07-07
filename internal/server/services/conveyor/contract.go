package conveyor

import "net/http"

type logger interface {
	Handle(h http.HandlerFunc) http.HandlerFunc
}
