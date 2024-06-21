package handlers

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGaugeFetcher(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name    string
		request string
		vars    map[string]string
		method  string
		fields  map[string]float64
		want    want
	}{
		{
			name:    "StatusNotFound test without metric",
			request: "/value/counter/",
			vars:    map[string]string{},
			method:  http.MethodGet,
			fields:  map[string]float64{},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain",
				body:        "",
			},
		},
		{
			name:    "StatusNotFound test with empty metric in storage",
			request: "/value/counter/",
			vars:    map[string]string{vars.Metric: "test"},
			method:  http.MethodGet,
			fields:  map[string]float64{},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain",
				body:        "",
			},
		},
		{
			name:    "StatusOK test",
			request: "/value/counter/test",
			vars:    map[string]string{vars.Metric: "test"},
			method:  http.MethodGet,
			fields:  map[string]float64{"test": 10.1},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain",
				body:        "10.1",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			rctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			for k, v := range test.vars {
				rctx.URLParams.Add(k, v)
			}

			storage := memstorage.New()
			for k, v := range test.fields {
				storage.SetGauge(k, v)
			}
			w := httptest.NewRecorder()
			GaugeFetcher(&storage)(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			body, err := io.ReadAll(result.Body)
			assert.Equal(t, []byte(test.want.body), body)
			assert.NoError(t, err)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
