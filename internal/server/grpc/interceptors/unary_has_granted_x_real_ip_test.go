package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestUnaryHasGrantedXRealIP(t *testing.T) {
	tests := []struct {
		name          string
		ip            string
		trustedSubnet string
		wantErr       bool
	}{
		{
			name:    "Test with empty trustedSubnet",
			wantErr: false,
		},
		{
			name:          "Test with empty ip",
			trustedSubnet: "192.168.1.0/24",
			wantErr:       false,
		},
		{
			name:          "Test with wrong subnet",
			trustedSubnet: "192.168.1.1.0/24",
			ip:            "193.168.1.1",
			wantErr:       false,
		},
		{
			name:          "Test with not granted ip",
			trustedSubnet: "192.168.1.0/24",
			ip:            "193.168.1.1",
			wantErr:       true,
		},
		{
			name:          "Test with granted ip",
			trustedSubnet: "192.168.1.0/24",
			ip:            "192.168.1.4",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UnaryHasGrantedXRealIP(tt.trustedSubnet)
			ctx := context.Background()
			md := metadata.New(map[string]string{"X-Real-IP": tt.ip})
			ctx = metadata.NewIncomingContext(ctx, md)
			handler := func(ctx context.Context, req any) (any, error) {
				return nil, nil
			}
			_, err := f(ctx, "", nil, handler)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
