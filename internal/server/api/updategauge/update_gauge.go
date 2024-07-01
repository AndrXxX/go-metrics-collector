package updategauge

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type updateGaugeHandler struct {
	u updater
}

func (h *updateGaugeHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric := chi.URLParam(r, vars.Metric)
	value := chi.URLParam(r, vars.Value)
	if err := h.u.Update(metric, value); err == nil {
		w.WriteHeader(http.StatusOK)
		return true
	}
	w.WriteHeader(http.StatusBadRequest)
	return false
}

func New(u updater) *updateGaugeHandler {
	return &updateGaugeHandler{u}
}
