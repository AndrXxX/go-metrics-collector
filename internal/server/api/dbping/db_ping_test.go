package dbping

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDBChecker struct {
	err error
}

func (c testDBChecker) Check(_ context.Context) error {
	return c.err
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		c    dbChecker
		want *dbPingHandler
	}{
		{
			name: "Test with checker",
			c:    testDBChecker{},
			want: &dbPingHandler{testDBChecker{}},
		},
		{
			name: "Test without checker",
			want: &dbPingHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(tt.c))
		})
	}
}

func Test_dbPingHandler_Handler(t *testing.T) {
	tests := []struct {
		name     string
		c        dbChecker
		wantCode int
	}{
		{
			name:     "StatusInternalServerError",
			c:        testDBChecker{err: fmt.Errorf("test error")},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "StatusOK",
			c:        testDBChecker{},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &dbPingHandler{c: tt.c}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			h.Handler()(w, r)
			result := w.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
