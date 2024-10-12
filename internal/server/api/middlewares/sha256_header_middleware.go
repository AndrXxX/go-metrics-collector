package middlewares

import (
	"bytes"
	"net/http"
)

type sha256HeaderMiddleware struct {
	hg  SHA256hashGenerator
	key string
}

func (m *sha256HeaderMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Handle(w, r, next)
	})
}

func (m *sha256HeaderMiddleware) Handle(w http.ResponseWriter, r *http.Request, next http.Handler) {
	w = m.processWriter(w)
	if next != nil {
		next.ServeHTTP(w, r)
	}
}

func (m *sha256HeaderMiddleware) processWriter(w http.ResponseWriter) http.ResponseWriter {
	if m.key == "" {
		return w
	}
	return &sha256RequestWriter{m.hg, w, &bytes.Buffer{}, m.key}
}

func AddSHA256HashHeader(hg SHA256hashGenerator, key string) *sha256HeaderMiddleware {
	return &sha256HeaderMiddleware{hg, key}
}

type sha256RequestWriter struct {
	hg   SHA256hashGenerator
	orig http.ResponseWriter
	buff *bytes.Buffer
	key  string
}

func (w *sha256RequestWriter) Header() http.Header {
	return w.orig.Header()
}

func (w *sha256RequestWriter) Write(data []byte) (int, error) {
	w.buff.Write(data)
	return w.orig.Write(data)
}

func (w *sha256RequestWriter) WriteHeader(statusCode int) {
	w.Header().Add("HashSHA256", w.hg.Generate(w.key, w.buff.Bytes()))
	w.orig.WriteHeader(statusCode)
}
