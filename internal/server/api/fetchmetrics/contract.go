package fetchmetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type storageProvider interface {
	GetStorage(name string) interfaces.MetricsStorage
}

type stringifier interface {
	String(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
