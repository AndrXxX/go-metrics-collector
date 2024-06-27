package updategauge

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type updateGaugeHandler struct {
	s guStorage
}

func (h *updateGaugeHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric := chi.URLParam(r, vars.Metric)
	value := chi.URLParam(r, vars.Value)
	if converted, err := strconv.ParseFloat(value, 64); err == nil {
		h.s.Insert(metric, converted)
		w.WriteHeader(http.StatusOK)
		return true
	}
	w.WriteHeader(http.StatusBadRequest)
	return false
}

func New(s guStorage) *updateGaugeHandler {
	return &updateGaugeHandler{s}
}
