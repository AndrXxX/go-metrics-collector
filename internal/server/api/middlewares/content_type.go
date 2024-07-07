package middlewares

import (
	"net/http"
)

type contentType struct {
	ct string
}

func (m *contentType) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", m.ct)
	if next != nil {
		next(w, r)
	}
}

func SetContentType(ct string) *contentType {
	return &contentType{ct}
}
