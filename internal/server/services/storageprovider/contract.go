package storageprovider

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type storage interface {
	Insert(metric string, value *models.Metrics)
	Get(metric string) (value *models.Metrics, ok bool)
}
