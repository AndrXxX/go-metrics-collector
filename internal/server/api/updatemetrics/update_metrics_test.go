package updatemetrics

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsidentifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricstringifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storageprovider"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricsHandlerGaugeHandle(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		vars    map[string]string
		method  string
		want    want
	}{
		{
			name:    "StatusOK test",
			request: "/update/gauge/test/10.1",
			vars:    map[string]string{vars.MetricType: metrics.Gauge, vars.Metric: "test", vars.Value: "10.1"},
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "StatusBadRequest test",
			request: "/update/gauge/test/aaa",
			vars:    map[string]string{vars.MetricType: metrics.Gauge, vars.Metric: "test", vars.Value: "aaa"},
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	storage := memory.New[*models.Metrics]()
	sp := storageprovider.New[interfaces.MetricsStorage]()
	sp.RegisterStorage(metrics.Gauge, &storage)
	h := New(metricsupdater.New(sp), metricstringifier.MetricsEmptyStringifier{}, metricsidentifier.NewURLIdentifier())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			ctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			for k, v := range test.vars {
				ctx.URLParams.Add(k, v)
			}

			w := httptest.NewRecorder()
			h.Handle(w, request, nil)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateMetricsHandlerCounterHandle(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		vars    map[string]string
		method  string
		want    want
	}{
		{
			name:    "StatusOK test",
			request: "/update/counter/test/10",
			vars:    map[string]string{vars.MetricType: metrics.Counter, vars.Metric: "test", vars.Value: "10"},
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "StatusBadRequest test",
			request: "/update/counter/test/aaa",
			vars:    map[string]string{vars.MetricType: metrics.Counter, vars.Metric: "test", vars.Value: "aaa"},
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	storage := memory.New[*models.Metrics]()
	sp := storageprovider.New[interfaces.MetricsStorage]()
	sp.RegisterStorage(metrics.Counter, &storage)
	h := New(metricsupdater.New(sp), metricstringifier.MetricsEmptyStringifier{}, metricsidentifier.NewURLIdentifier())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			ctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			for k, v := range test.vars {
				ctx.URLParams.Add(k, v)
			}

			w := httptest.NewRecorder()

			h.Handle(w, request, nil)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
