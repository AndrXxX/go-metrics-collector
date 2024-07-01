package metricsupdater

import "github.com/AndrXxX/go-metrics-collector/internal/server/models"

type metricsSetter interface {
	Set(m *models.Metrics, value string) error
}
