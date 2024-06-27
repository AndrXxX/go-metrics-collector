package update_gauge

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type guStorage interface {
	Insert(metric string, value float64)
}

func GaugeUpdater(s guStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.Metric)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value := chi.URLParam(r, vars.Value)
		if converted, err := strconv.ParseFloat(value, 64); err == nil {
			s.Insert(metric, converted)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
