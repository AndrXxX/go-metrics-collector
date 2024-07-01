package metricsidentifier

import (
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type jsonMetricsIdentifier struct {
}

func (i *jsonMetricsIdentifier) Process(r *http.Request) (*models.Metrics, error) {
	var m *models.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	return m, err
}

func NewJsonIdentifier() *jsonMetricsIdentifier {
	return &jsonMetricsIdentifier{}
}
