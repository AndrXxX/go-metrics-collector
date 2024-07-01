package metricsidentifier

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type urlMetricsIdentifier struct {
	mType string
}

func (i *urlMetricsIdentifier) Process(r *http.Request) (*models.Metrics, error) {
	return &models.Metrics{
		ID:    chi.URLParam(r, vars.Metric),
		MType: i.mType,
	}, nil
}

func NewURLIdentifier(mType string) *urlMetricsIdentifier {
	return &urlMetricsIdentifier{mType}
}
