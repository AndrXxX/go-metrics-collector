package middlewares

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func HasMetricOr404() interfaces.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (ok bool) {
		metric := chi.URLParam(r, vars.Metric)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return false
		}
		return true
	}
}
