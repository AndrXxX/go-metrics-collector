package middlewares

import (
	"bytes"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/server/services/sha256"
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
	return &sha256.RequestWriter{HG: m.hg, OriginalWriter: w, Buffer: &bytes.Buffer{}, Key: m.key}
}

// AddSHA256HashHeader возвращает middleware, которая добавляет хеш по ключу HashSHA256
func AddSHA256HashHeader(hg hashGenerator, key string) *sha256HeaderMiddleware {
	return &sha256HeaderMiddleware{hg, key}
}
