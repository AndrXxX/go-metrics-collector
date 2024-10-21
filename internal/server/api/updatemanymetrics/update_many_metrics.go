package updatemanymetrics

import (
	"encoding/json"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type updateManyMetricsHandler struct {
	u updater
}

// Handler возвращает http.HandlerFunc для обновления нескольких метрик из запроса
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

// New возвращает обработчик updateManyMetricsHandler для обновления нескольких метрик
func New(u updater) *updateManyMetricsHandler {
	return &updateManyMetricsHandler{u}
}
