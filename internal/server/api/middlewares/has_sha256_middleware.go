package middlewares

import (
	"bytes"
	"io"
	"net/http"
)

type hasCorrectSHA256HashMiddleware struct {
	hg SHA256hashGenerator
}

func (m *hasCorrectSHA256HashMiddleware) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !m.check(r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if next != nil {
		next(w, r)
	}
}

func (m *hasCorrectSHA256HashMiddleware) check(r *http.Request) bool {
	if m.hg == nil {
		return true
	}
	requestHash := r.Header.Get("HashSHA256")
	if requestHash == "" {
		return true
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	actualHash := m.hg.Generate(data)
	return actualHash == requestHash
}

func HasCorrectSHA256HashOr500(hg SHA256hashGenerator) *hasCorrectSHA256HashMiddleware {
	return &hasCorrectSHA256HashMiddleware{hg}
}
