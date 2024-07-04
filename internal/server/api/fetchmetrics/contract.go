package fetchmetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type storage[T any] interface {
	Insert(metric string, value T)
	Get(metric string) (value T, ok bool)
}

type formatter interface {
	Format(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
