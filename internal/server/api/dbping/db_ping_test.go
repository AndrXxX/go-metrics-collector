package dbping

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDbChecker struct {
	err error
}

func (c testDbChecker) Check(_ context.Context) error {
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
			c:    testDbChecker{},
			want: &dbPingHandler{testDbChecker{}},
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
			c:        testDbChecker{err: fmt.Errorf("test error")},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "StatusOK",
			c:        testDbChecker{},
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
		})
	}
}
