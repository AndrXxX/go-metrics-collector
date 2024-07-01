package interfaces

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsStorage interface {
	Insert(metric string, value *models.Metrics)
	Get(metric string) (value *models.Metrics, ok bool)
	All() map[string]*models.Metrics
}
