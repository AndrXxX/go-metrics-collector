package agent

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	grpsserv "github.com/AndrXxX/go-metrics-collector/internal/agent/services/grpc"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/grpc/dealoptions"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
)

func WithGRPCMetricsUploader(hg hashGenerator) Option {
	return func(a *agent) {
		if a.c.Common.GRPCHost == "" {
			return
		}
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
			grpc.WithUnaryInterceptor(dealoptions.WithXRealIP(a.c.Common.Host)),
			grpc.WithUnaryInterceptor(dealoptions.WithSHA256(hg, a.c.Common.Key)),
		}
		updater := grpsserv.NewGRPCMetricsUpdater(a.c.Common.GRPCHost, opts)
		a.processors.Add(metricsuploader.NewGRPCUploader(updater))
	}
}
