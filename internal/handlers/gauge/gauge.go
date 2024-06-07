package gauge

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func UpdateHandler(s repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.METRIC)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value := chi.URLParam(r, vars.VALUE)
		if converted, err := strconv.ParseFloat(value, 64); err == nil {
			s.SetGauge(metric, converted)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func FetchHandler(s repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metric := chi.URLParam(r, vars.METRIC)
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if val, err := s.GetGauge(metric); err == nil {
			_, _ = w.Write([]byte(fmt.Sprintf("%v", val)))
			w.WriteHeader(http.StatusOK)
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
