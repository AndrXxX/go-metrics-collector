package dealoptions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestWithXRealIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantIP  string
		wantErr bool
	}{
		{
			name:   "Test with empty ip",
			wantIP: "",
		},
		{
			name:   "Test with ip 192.168.0.12",
			ip:     "192.168.0.12",
			wantIP: "192.168.0.12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := WithXRealIP(tt.ip)
			ctx := context.Background()
			var ip string
			invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				if md, ok := metadata.FromOutgoingContext(ctx); ok {
					values := md.Get("X-Real-IP")
					if len(values) > 0 {
						ip = values[0]
					}
				}
				return nil
			}
			err := f(ctx, "", "", nil, nil, invoker)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantIP, ip)
		})
	}
}
