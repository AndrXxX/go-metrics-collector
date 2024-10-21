package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
)

type hasMetricOr404 struct {
}

// Handler возвращает http.HandlerFunc
func (m *hasMetricOr404) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metric := chi.URLParam(r, vars.Metric)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// HasMetricOr404 возвращает middleware, которая проверяет наличие названия метрики в запросе или возвращает 404
func HasMetricOr404() *hasMetricOr404 {
	return &hasMetricOr404{}
}
