package dealoptions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type testHg struct {
	mock.Mock
}

func (g *testHg) Generate(key string, data []byte) string {
	args := g.Called(key, data)
	return args.String(0)
}

func TestWithSHA256(t *testing.T) {
	tests := []struct {
		name     string
		hg       hashGenerator
		key      string
		request  any
		wantHash string
		wantErr  bool
	}{
		{
			name:    "Test with empty key",
			wantErr: false,
		},
		{
			name:    "Test with error on marshall",
			key:     "test",
			request: func() {},
			wantErr: true,
		},
		{
			name: "Test with hash `hashedResult`",
			hg: func() hashGenerator {
				hg := testHg{}
				hg.On("Generate", mock.Anything, mock.Anything).Return("hashedResult")
				return &hg
			}(),
			key:      "test",
			wantHash: "hashedResult",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := WithSHA256(tt.hg, tt.key)
			ctx := context.Background()
			var hash string
			invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				if md, ok := metadata.FromOutgoingContext(ctx); ok {
					values := md.Get("HashSHA256")
					if len(values) > 0 {
						hash = values[0]
					}
				}
				return nil
			}
			err := f(ctx, "", tt.request, nil, nil, invoker)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantHash, hash)
		})
	}
}
