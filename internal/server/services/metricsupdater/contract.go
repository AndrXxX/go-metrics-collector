package metricsupdater

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type metricsSetter interface {
	Set(m *models.Metrics, value string) error
}

type storage interface {
	Insert(metric string, value *models.Metrics)
	Get(metric string) (value *models.Metrics, ok bool)
}
