package metricsidentifier

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsvaluesetter"
)

type urlMetricsIdentifier struct {
}

func (i *urlMetricsIdentifier) Process(r *http.Request) (*models.Metrics, error) {
	m := models.Metrics{
		ID:    chi.URLParam(r, vars.Metric),
		MType: chi.URLParam(r, vars.MetricType),
	}
	value := chi.URLParam(r, vars.Value)
	if value == "" {
		return &m, nil
	}
	setter, ok := metricsvaluesetter.Factory().SetterByType(m.MType)
	if !ok {
		return nil, fmt.Errorf("unknown metrics type %q", m.MType)
	}
	err := setter.Set(&m, value)
	return &m, err
}

func NewURLIdentifier() *urlMetricsIdentifier {
	return &urlMetricsIdentifier{}
}
