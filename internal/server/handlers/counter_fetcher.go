package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func CounterFetcher(s repositories.Storage[int64]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.Metric)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		val, ok := s.Get(metric)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(fmt.Sprintf("%d", val)))
		if err != nil {
			logger.Log.Error("Failed to write counter response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
