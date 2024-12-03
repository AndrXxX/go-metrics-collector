package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		res, err := handler(ctx, req)

		duration := time.Since(start)
		var code codes.Code
		if err != nil {
			if e, ok := status.FromError(err); ok {
				code = e.Code()
			}
		}

		logger.Log.Info(
			"got incoming gRPC request",
			zap.String("method", info.FullMethod),
			zap.Uint32("code", uint32(code)),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return res, err
	}
}
