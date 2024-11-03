package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
)

func TestHasMetricOr404(t *testing.T) {
	tests := []struct {
		name string
		want *hasMetricOr404
	}{
		{
			name: "Test OK",
			want: &hasMetricOr404{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, HasMetricOr404())
		})
	}
}

func Test_hasMetricOr404_Handler(t *testing.T) {
	tests := []struct {
		name     string
		vars     map[string]string
		next     http.Handler
		wantCode int
	}{
		{
			name:     "Test if metric not exist",
			vars:     map[string]string{},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test if metric exist and has next func",
			vars: map[string]string{vars.Metric: "test"},
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &hasMetricOr404{}
			h := m.Handler(tt.next)
			request := httptest.NewRequest(http.MethodGet, "/value/counter/test", nil)
			rc := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rc))
			for k, v := range tt.vars {
				rc.URLParams.Add(k, v)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
