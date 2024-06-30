package fetchcounter

import "github.com/AndrXxX/go-metrics-collector/internal/server/models"

type storage interface {
	Get(metric string) (value *models.Metrics, ok bool)
}
