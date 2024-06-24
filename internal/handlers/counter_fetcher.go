package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func CounterFetcher(s repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.METRIC)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if val, err := s.GetCounter(metric); err == nil {
			_, _ = w.Write([]byte(fmt.Sprintf("%d", val)))
			w.WriteHeader(http.StatusOK)
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
