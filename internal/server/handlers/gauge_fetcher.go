package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GaugeFetcher(s repositories.Storage[float64]) http.HandlerFunc {
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
