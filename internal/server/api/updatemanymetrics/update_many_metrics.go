package updatemanymetrics

import (
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type updateManyMetricsHandler struct {
	u updater
}

// Handler returns HandlerFunc to update metrics from request
func (h *updateManyMetricsHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// New Return updateManyMetricsHandler
func New(u updater) *updateManyMetricsHandler {
	return &updateManyMetricsHandler{u}
}
