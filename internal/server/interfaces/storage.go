package interfaces

import (
	"context"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type MetricsStorage interface {
	Insert(ctx context.Context, metric string, value *models.Metrics)
	Get(ctx context.Context, metric string) (value *models.Metrics, ok bool)
	All(ctx context.Context) map[string]*models.Metrics
	Delete(ctx context.Context, metric string) (ok bool)
}
