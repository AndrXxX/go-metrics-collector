package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type requestLogger struct {
}

// Handler возвращает http.HandlerFunc
func (l *requestLogger) Handler(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info(
			"got incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	}
	return http.HandlerFunc(logFn)
}

// New возвращает обработчик для логирования запросов
func New() *requestLogger {
	return &requestLogger{}
}
