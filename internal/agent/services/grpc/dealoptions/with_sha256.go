package dealoptions

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func WithSHA256(hg hashGenerator, key string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if key == "" {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		encoded, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error on json marshall grapc request: %w", err)
		}
		md := metadata.New(map[string]string{"HashSHA256": hg.Generate(key, encoded)})
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
