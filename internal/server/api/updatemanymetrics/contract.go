package updatemanymetrics

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type updater interface {
	UpdateMany(context.Context, []*models.Metrics) error
}
