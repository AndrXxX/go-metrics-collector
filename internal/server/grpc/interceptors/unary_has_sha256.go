package interceptors

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func UnaryHasCorrectSHA256(hg hashGenerator, key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if key == "" {
			return handler(ctx, req)
		}
		var requestHash string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("HashSHA256")
			if len(values) > 0 {
				requestHash = values[0]
			}
		}
		if requestHash == "" {
			return handler(ctx, req)
		}
		encoded, err := json.Marshal(req)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("failed to marshal request: %v", err))
			return handler(ctx, req)
		}
		if requestHash != hg.Generate(key, encoded) {
			return nil, status.Errorf(codes.PermissionDenied, "Permission Denied")
		}
		return handler(ctx, req)
	}
}
