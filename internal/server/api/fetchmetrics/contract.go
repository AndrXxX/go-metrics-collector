package fetchmetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type storage interface {
	Get(metric string) (value *models.Metrics, ok bool)
}

type stringifier interface {
	String(m *models.Metrics) (string, error)
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
