package handlers

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type updater interface {
	Update(name string, value int64)
}

func CounterUpdater(u updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.Metric)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value := chi.URLParam(r, vars.Value)
		if converted, err := strconv.ParseInt(value, 10, 64); err == nil {
			u.Update(metric, converted)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
