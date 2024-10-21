package middlewares

import (
	"bytes"
	"net/http"
)

type sha256HeaderMiddleware struct {
	hg  hashGenerator
	key string
}

// Handler возвращает http.HandlerFunc
func (m *sha256HeaderMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = m.processWriter(w)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *sha256HeaderMiddleware) processWriter(w http.ResponseWriter) http.ResponseWriter {
	if m.key == "" {
		return w
	}
	return &sha256RequestWriter{m.hg, w, &bytes.Buffer{}, m.key}
}

// AddSHA256HashHeader возвращает middleware, которая добавляет хеш по ключу HashSHA256
func AddSHA256HashHeader(hg hashGenerator, key string) *sha256HeaderMiddleware {
	return &sha256HeaderMiddleware{hg, key}
}

type sha256RequestWriter struct {
	hg   hashGenerator
	orig http.ResponseWriter
	buff *bytes.Buffer
	key  string
}

// Header is implementation of http.ResponseWriter
func (w *sha256RequestWriter) Header() http.Header {
	return w.orig.Header()
}

// Write is implementation of http.ResponseWriter
func (w *sha256RequestWriter) Write(data []byte) (int, error) {
	w.buff.Write(data)
	return w.orig.Write(data)
}

// WriteHeader is implementation of http.ResponseWriter
func (w *sha256RequestWriter) WriteHeader(statusCode int) {
	w.Header().Add("HashSHA256", w.hg.Generate(w.key, w.buff.Bytes()))
	w.orig.WriteHeader(statusCode)
}
