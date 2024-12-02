package metricsuploader

import (
	"context"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	pb "github.com/AndrXxX/go-metrics-collector/internal/proto"
)

type urlBuilder interface {
	Build(params types.URLParams) string
}

type requestSender interface {
	Post(url string, contentType string, data []byte) error
}

type grpcMetricsUpdater interface {
	Update(ctx context.Context, list []*pb.Metric) error
}
