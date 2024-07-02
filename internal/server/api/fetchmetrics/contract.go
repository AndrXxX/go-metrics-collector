package fetchmetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type storage[T any] interface {
	Insert(metric string, value T)
	Get(metric string) (value T, ok bool)
}

type storageProvider[T any] interface {
	GetStorage(name string) T
}

type stringifier interface {
	String(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
