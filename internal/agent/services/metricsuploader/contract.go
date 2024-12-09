package metricsuploader

import (
	"context"

	mp "github.com/AndrXxX/go-metrics-collector/pkg/metricsproto"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
)

type urlBuilder interface {
	Build(params types.URLParams) string
}

type requestSender interface {
	Post(url string, contentType string, data []byte) error
}

type grpcMetricsUpdater interface {
	Update(ctx context.Context, list []*mp.Metric) error
}
