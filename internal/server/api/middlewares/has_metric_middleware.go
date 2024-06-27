package middlewares

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type hasMetricOr404 struct {
}

func (m *hasMetricOr404) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric := chi.URLParam(r, vars.Metric)
	if metric == "" {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	return true
}

func HasMetricOr404() *hasMetricOr404 {
	return &hasMetricOr404{}
}
