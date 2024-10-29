package middlewares

import (
	"net/http"
)

type contentType struct {
	ct string
}

// Handler возвращает http.HandlerFunc
func (m *contentType) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", m.ct)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// SetContentType возвращает middleware для установки Content-Type
func SetContentType(ct string) *contentType {
	return &contentType{ct}
}
