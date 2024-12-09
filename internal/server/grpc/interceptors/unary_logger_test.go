package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testLogger struct {
	msg    string
	fields map[string]zap.Field
}

func (l *testLogger) Info(msg string, fields ...zap.Field) {
	l.msg = msg
	l.fields = map[string]zap.Field{}
	for _, f := range fields {
		l.fields[f.Key] = f
	}
}

func TestUnaryLogger(t *testing.T) {
	type want struct {
		msg    string
		fields []string
		values map[string]any
	}
	tests := []struct {
		name   string
		err    error
		method string
		want   want
	}{
		{
			name:   "Test without error",
			method: "testMethod",
			want: want{
				msg:    "got incoming gRPC request",
				fields: []string{"method", "code", "duration"},
				values: map[string]any{
					"method": "testMethod",
				},
			},
		},
		{
			name:   "Test with error",
			method: "testMethod",
			err:    status.Errorf(codes.PermissionDenied, "Permission Denied"),
			want: want{
				msg:    "got incoming gRPC request",
				fields: []string{"method", "code", "duration", "error"},
				values: map[string]any{
					"method": "testMethod",
					"code":   int64(codes.PermissionDenied),
					"error":  status.Errorf(codes.PermissionDenied, "Permission Denied"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := testLogger{}
			f := UnaryLogger(&l)
			ctx := context.Background()
			handler := func(ctx context.Context, req any) (any, error) {
				return nil, tt.err
			}
			_, err := f(ctx, nil, &grpc.UnaryServerInfo{FullMethod: tt.method}, handler)
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.want.msg, l.msg)
			for _, field := range tt.want.fields {
				assert.NotEmpty(t, l.fields[field])
			}
			for f, v := range tt.want.values {
				if l.fields[f].String != "" {
					assert.Equal(t, v, l.fields[f].String)
					continue
				}
				if l.fields[f].Integer > 0 {
					assert.Equal(t, v, l.fields[f].Integer)
					continue
				}
				if err, ok := l.fields[f].Interface.(error); ok {
					assert.Equal(t, v, err)
					continue
				}
			}
		})
	}
}
