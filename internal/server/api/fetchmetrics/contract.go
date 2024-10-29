package fetchmetrics

import (
	"context"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type storage[T any] interface {
	Insert(ctx context.Context, metric string, value T)
	Get(ctx context.Context, metric string) (value T, ok bool)
}

type formatter interface {
	Format(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}

type metricsChecker interface {
	IsValid(m *models.Metrics) bool
}
