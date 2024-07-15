package interfaces

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsStorage interface {
	Insert(metric string, value *models.Metrics)
	Get(metric string) (value *models.Metrics, ok bool)
	All(ctx context.Context) map[string]*models.Metrics
}
