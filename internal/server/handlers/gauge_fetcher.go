package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type gfStorage interface {
	Get(metric string) (value float64, ok bool)
}

func GaugeFetcher(s gfStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.Metric)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if val, ok := s.Get(metric); ok {
			_, _ = fmt.Fprintf(w, "%v", val)
			w.WriteHeader(http.StatusOK)
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
