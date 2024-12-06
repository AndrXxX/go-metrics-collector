package dbchecker

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConn struct {
	err error
}

func (c testConn) PingContext(_ context.Context) error {
	return c.err
}

func Test_dbChecker_Check(t *testing.T) {
	tests := []struct {
		name    string
		db      conn
		wantErr bool
	}{
		{
			name:    "Test with nil connection",
			db:      nil,
			wantErr: true,
		},
		{
			name:    "Test with error on ping connection",
			db:      &testConn{err: errors.New("test")},
			wantErr: true,
		},
		{
			name:    "Test OK",
			db:      &testConn{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.db)
			require.Equal(t, tt.wantErr, c.Check(context.Background()) != nil)
		})
	}
}
