package interceptors

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func UnaryHasGrantedXRealIP(trustedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if trustedSubnet == "" {
			return handler(ctx, req)
		}
		var ip string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("X-Real-Ip")
			if len(values) > 0 {
				ip = values[0]
			}
		}
		if ip == "" {
			return handler(ctx, req)
		}
		_, ipNet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("TrustedSubnet %s is not valid", trustedSubnet))
			return handler(ctx, req)
		}
		if !ipNet.Contains(net.ParseIP(ip)) {
			return nil, status.Errorf(codes.PermissionDenied, "Permission Denied")
		}
		return handler(ctx, req)
	}
}
