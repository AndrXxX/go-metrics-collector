package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type testHg struct {
	mock.Mock
}

func (g *testHg) Generate(key string, data []byte) string {
	args := g.Called(key, data)
	return args.String(0)
}

func TestUnaryHasCorrectSHA256(t *testing.T) {
	tests := []struct {
		name    string
		hg      hashGenerator
		key     string
		hash    string
		req     interface{}
		wantErr bool
	}{
		{
			name:    "Test with empty key",
			wantErr: false,
		},
		{
			name:    "Test with empty hash",
			key:     "testKey",
			wantErr: false,
		},
		{
			name:    "Test with error on marshal request",
			key:     "testKey",
			hash:    "hashedResult",
			req:     func() {},
			wantErr: false,
		},
		{
			name: "Test with wrong hash",
			key:  "testKey",
			hash: "hashedResultWrong",
			hg: func() hashGenerator {
				hg := testHg{}
				hg.On("Generate", mock.Anything, mock.Anything).Return("hashedResult")
				return &hg
			}(),
			wantErr: true,
		},
		{
			name: "Test with correct hash",
			key:  "testKey",
			hash: "hashedResult",
			hg: func() hashGenerator {
				hg := testHg{}
				hg.On("Generate", mock.Anything, mock.Anything).Return("hashedResult")
				return &hg
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UnaryHasCorrectSHA256(tt.hg, tt.key)
			ctx := context.Background()
			md := metadata.New(map[string]string{"HashSHA256": tt.hash})
			ctx = metadata.NewIncomingContext(ctx, md)
			handler := func(ctx context.Context, req any) (any, error) {
				return nil, nil
			}
			_, err := f(ctx, tt.req, nil, handler)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
