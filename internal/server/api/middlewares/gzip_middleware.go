package middlewares

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/server/services/gzipcompressor"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type gzipMiddleware struct {
}

// Handler возвращает http.HandlerFunc
func (m *gzipMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := gzipcompressor.NewCompressWriter(w)
			ow = cw
			defer func() {
				_ = cw.Close()
			}()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := gzipcompressor.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Error("Error creating gzip compressor", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				_ = cr.Close()
			}()
		}
		if next != nil {
			next.ServeHTTP(ow, r)
		}
	})
}

// CompressGzip возвращает middleware для сжатия ответа и распаковки запроса
func CompressGzip() *gzipMiddleware {
	return &gzipMiddleware{}
}
