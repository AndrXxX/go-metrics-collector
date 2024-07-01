package metricsidentifier

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type urlMetricsIdentifier struct {
}

func (i *urlMetricsIdentifier) Process(r *http.Request) (*models.Metrics, error) {
	return &models.Metrics{
		ID:    chi.URLParam(r, vars.Metric),
		MType: chi.URLParam(r, vars.MetricType),
	}, nil
}

func NewURLIdentifier() *urlMetricsIdentifier {
	return &urlMetricsIdentifier{}
}
