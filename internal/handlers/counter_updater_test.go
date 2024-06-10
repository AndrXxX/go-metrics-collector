package handlers

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounterUpdater(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		vars    map[string]string
		method  string
		want    want
	}{
		{
			name:    "StatusNotFound test",
			request: "/update/counter/",
			vars:    map[string]string{},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "StatusOK test",
			request: "/update/counter/test/10",
			vars:    map[string]string{vars.Metric: "test", vars.Value: "10"},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "StatusBadRequest test",
			request: "/update/counter/test/dsff",
			vars:    map[string]string{vars.Metric: "test", vars.Value: "dsff"},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	storage := memstorage.New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			rctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			for k, v := range test.vars {
				rctx.URLParams.Add(k, v)
			}

			w := httptest.NewRecorder()
			CounterUpdater(&storage)(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
