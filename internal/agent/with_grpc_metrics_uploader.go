package agent

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	grpsserv "github.com/AndrXxX/go-metrics-collector/internal/agent/services/grpc"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/grpc/dealoptions"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func WithGRPCMetricsUploader(hg hashGenerator, tlsProvider tlsConfigProvider) Option {
	return func(a *agent) {
		if a.c.Common.GRPCHost == "" {
			return
		}
		var opts types.ItemsList[grpc.DialOption]

		tlsConf, err := tlsProvider.Fetch()
		if err != nil {
			logger.Log.Error("failed to fetch tlsConf", zap.Error(err))
		}
		if tlsConf != nil {
			opts.Add(grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
		} else {
			opts.Add(grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
		opts.Add(grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
		opts.Add(grpc.WithUnaryInterceptor(dealoptions.WithXRealIP(a.c.Common.Host)))
		opts.Add(grpc.WithUnaryInterceptor(dealoptions.WithSHA256(hg, a.c.Common.Key)))
		updater := grpsserv.NewGRPCMetricsUpdater(a.c.Common.GRPCHost, opts)
		a.processors.Add(metricsuploader.NewGRPCUploader(updater))
	}
}
