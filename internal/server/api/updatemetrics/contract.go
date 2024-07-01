package updatemetrics

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"net/http"
)

type updater interface {
	Update(m *models.Metrics, value string) error
}

type identifier interface {
	Process(r *http.Request) (*models.Metrics, error)
}
