package fetchmetrics

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricstringifier"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchMetricsHandlerGaugeHandle(t *testing.T) {
	type want struct {
		statusCode int
		body       string
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
			name:    "StatusNotFound test with empty metric in storage",
			request: "/value/counter/",
			vars:    map[string]string{vars.Metric: "test"},
			method:  http.MethodGet,
			fields:  map[string]float64{},
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:    "StatusOK test",
			request: "/value/counter/test",
			vars:    map[string]string{vars.Metric: "test"},
			method:  http.MethodGet,
			fields:  map[string]float64{"test": 10.1},
			want: want{
				statusCode: http.StatusOK,
				body:       "10.1",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			ctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			for k, v := range test.vars {
				ctx.URLParams.Add(k, v)
			}

			storage := memory.New[*models.Metrics]()
			for k, v := range test.fields {
				storage.Insert(k, &models.Metrics{
					ID:    k,
					MType: metrics.Gauge,
					Value: &v,
				})
			}
			w := httptest.NewRecorder()
			h := New(&storage, metricstringifier.MetricsValueStringifier{})
			h.Handle(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
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

func TestFetchMetricsHandlerCounterHandle(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		request string
		vars    map[string]string
		method  string
		fields  map[string]int64
		want    want
	}{
		{
			name:    "StatusNotFound test with empty metric in storage",
			request: "/value/counter/",
			vars:    map[string]string{vars.Metric: "test"},
			method:  http.MethodGet,
			fields:  map[string]int64{},
			want: want{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name:    "StatusOK test",
			request: "/value/counter/test",
			vars:    map[string]string{vars.Metric: "test"},
			method:  http.MethodGet,
			fields:  map[string]int64{"test": 10},
			want: want{
				statusCode: http.StatusOK,
				body:       "10",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			ctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			for k, v := range test.vars {
				ctx.URLParams.Add(k, v)
			}

			storage := memory.New[*models.Metrics]()
			for k, v := range test.fields {
				storage.Insert(k, &models.Metrics{
					ID:    k,
					MType: metrics.Counter,
					Delta: &v,
				})
			}
			w := httptest.NewRecorder()
			h := New(&storage, metricstringifier.MetricsValueStringifier{})
			h.Handle(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
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
