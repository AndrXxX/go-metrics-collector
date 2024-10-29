package updatemetrics

import (
	"context"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type updater interface {
	Update(ctx context.Context, m *models.Metrics) (*models.Metrics, error)
}

type formatter interface {
	Format(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
