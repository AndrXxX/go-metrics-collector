package metricsidentifier

import (
	"encoding/json"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type jsonMetricsIdentifier struct {
}

func (i *jsonMetricsIdentifier) Process(r *http.Request) (*models.Metrics, error) {
	var m *models.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	return m, err
}

func NewJSONIdentifier() *jsonMetricsIdentifier {
	return &jsonMetricsIdentifier{}
}
