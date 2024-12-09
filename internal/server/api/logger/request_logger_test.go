package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_requestLogger_Handler(t *testing.T) {
	tests := []struct {
		name     string
		next     http.Handler
		wantCode int
	}{
		{
			name: "test StatusOK",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantCode: http.StatusOK,
		},
		{
			name: "test StatusForbidden",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}),
			wantCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New().Handler(tt.next)
			request := httptest.NewRequest("", "/test", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
