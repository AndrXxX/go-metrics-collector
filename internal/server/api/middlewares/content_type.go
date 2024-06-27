package middlewares

import (
	"net/http"
)

type contentType struct {
	ct string
}

func (m *contentType) Handle(w http.ResponseWriter, _ *http.Request) (ok bool) {
	w.Header().Set("Content-Type", m.ct)
	return true
}

func SetContentType(ct string) *contentType {
	return &contentType{ct}
}
