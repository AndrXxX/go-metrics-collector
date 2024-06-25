package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type cfStorage interface {
	Get(metric string) (value int64, ok bool)
}

func CounterFetcher(s cfStorage) http.HandlerFunc {
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
