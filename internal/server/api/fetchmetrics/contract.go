package fetchmetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type storageProvider[T any] interface {
	GetStorage(name string) T
}

type stringifier interface {
	String(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
