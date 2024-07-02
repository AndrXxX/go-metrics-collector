package middlewares

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type hasMetricOr404 struct {
}

func (m *hasMetricOr404) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	metric := chi.URLParam(r, vars.Metric)
	if metric == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if next != nil {
		next(w, r)
	}
}

func HasMetricOr404() *hasMetricOr404 {
	return &hasMetricOr404{}
}
