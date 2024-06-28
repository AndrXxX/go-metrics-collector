package middlewares

import (
	"net/http"
)

type httpError struct {
	code int
}

func (m *httpError) Handle(w http.ResponseWriter, _ *http.Request) (ok bool) {
	http.Error(w, http.StatusText(m.code), m.code)
	return true
}

func SetHTTPError(code int) *httpError {
	return &httpError{code}
}
