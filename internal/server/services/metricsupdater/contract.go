package metricsupdater

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
)

type storageProvider interface {
	GetStorage(name string) interfaces.MetricsStorage
}
