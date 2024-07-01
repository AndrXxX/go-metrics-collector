package updatemetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type updateMetricsHandler struct {
	u updater
	i identifier
}

func (h *updateMetricsHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric, err := h.i.Process(r)
	if metric == nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	value := chi.URLParam(r, vars.Value)
	if err := h.u.Update(metric, value); err == nil {
		w.WriteHeader(http.StatusOK)
		return true
	}
	w.WriteHeader(http.StatusBadRequest)
	return false
}

func New(u updater, i identifier) *updateMetricsHandler {
	return &updateMetricsHandler{u, i}
}
