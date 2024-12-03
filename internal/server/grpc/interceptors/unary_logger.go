package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		res, err := handler(ctx, req)

		duration := time.Since(start)

		logger.Log.Info(
			"got incoming gRPC request",
			zap.String("method", info.FullMethod),
			//zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			//zap.Int("size", responseData.size),
		)
		return res, err
	}
}
