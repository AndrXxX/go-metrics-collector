package updatemanymetrics

import (
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type updateManyMetricsHandler struct {
	u updater
}

// Handle updates metrics from request
func (h *updateManyMetricsHandler) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var list []models.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.u.UpdateMany(r.Context(), list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	if next != nil {
		next(w, r)
	}
}

// New Return updateManyMetricsHandler
func New(u updater) *updateManyMetricsHandler {
	return &updateManyMetricsHandler{u}
}
