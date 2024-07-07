package middlewares

import (
	"net/http"
)

type httpError struct {
	code int
}

func (m *httpError) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	http.Error(w, http.StatusText(m.code), m.code)
	if next != nil {
		next(w, r)
	}
}

func SetHTTPError(code int) *httpError {
	return &httpError{code}
}
