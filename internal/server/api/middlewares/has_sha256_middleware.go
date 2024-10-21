package middlewares

import (
	"bytes"
	"io"
	"net/http"
)

type hasCorrectSHA256HashMiddleware struct {
	hg  hashGenerator
	key string
}

// Handler возвращает http.HandlerFunc
func (m *hasCorrectSHA256HashMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.check(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *hasCorrectSHA256HashMiddleware) check(r *http.Request) bool {
	if m.key == "" {
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

	actualHash := m.hg.Generate(m.key, data)
	return actualHash == requestHash
}

// HasCorrectSHA256HashOr500 возвращает middleware, которая проверяет наличие и правильность хеша в запросе по ключу HashSHA256
func HasCorrectSHA256HashOr500(hg hashGenerator, key string) *hasCorrectSHA256HashMiddleware {
	return &hasCorrectSHA256HashMiddleware{hg, key}
}
