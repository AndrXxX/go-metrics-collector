package fetchmetrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricschecker"
)

type testIdentifier struct {
	err error
	m   models.Metrics
}

func (i testIdentifier) Process(_ *http.Request) (*models.Metrics, error) {
	return &i.m, i.err
}

type testFormatter struct {
	err error
	v   string
}

// Format возвращает значение метрики в виде строки
func (s testFormatter) Format(_ *models.Metrics) (string, error) {
	return s.v, s.err
}

func TestFetchMetricsHandlerGaugeHandle(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		fields map[string]float64
		i      identifier
		f      formatter
		want   want
	}{
		{
			name: "StatusNotFound test with not processed metric",
			i:    testIdentifier{err: fmt.Errorf("test error")},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "StatusBadRequest test with unknown metric type",
			i:    testIdentifier{m: models.Metrics{ID: "test", MType: "unknown"}},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "StatusNotFound test with empty metric in storage",
			i:    testIdentifier{m: models.Metrics{ID: "test", MType: metrics.Gauge}},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:   "StatusOK test",
			fields: map[string]float64{"test": 10.1},
			i:      testIdentifier{m: models.Metrics{ID: "test", MType: metrics.Gauge}},
			f:      testFormatter{v: "10.1"},
			want: want{
				statusCode: http.StatusOK,
				body:       "10.1",
			},
		},
		{
			name:   "StatusInternalServerError on format value",
			fields: map[string]float64{"test": 10.1},
			i:      testIdentifier{m: models.Metrics{ID: "test", MType: metrics.Gauge}},
			f:      testFormatter{err: fmt.Errorf("test error")},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/test", nil)

			storage := memory.New[*models.Metrics]()
			ctx := context.TODO()
			for k, v := range test.fields {
				storage.Insert(ctx, k, &models.Metrics{
					ID:    k,
					MType: metrics.Gauge,
					Value: &v,
				})
			}
			w := httptest.NewRecorder()
			h := New(&storage, test.f, test.i, metricschecker.New())
			h.Handler()(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			body, err := io.ReadAll(result.Body)
			assert.Equal(t, []byte(test.want.body), body)
			require.NoError(t, err)
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
		name   string
		fields map[string]int64
		i      identifier
		f      formatter
		want   want
	}{
		{
			name: "StatusBadRequest test with unknown metric type",
			i:    testIdentifier{m: models.Metrics{ID: "test", MType: "unknown"}},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "StatusNotFound test with empty metric in storage",
			i:    testIdentifier{m: models.Metrics{ID: "test", MType: metrics.Counter}},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:   "StatusOK test",
			fields: map[string]int64{"test": 10},
			i:      testIdentifier{m: models.Metrics{ID: "test", MType: metrics.Counter}},
			f:      testFormatter{v: "10"},
			want: want{
				statusCode: http.StatusOK,
				body:       "10",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/test", nil)

			ctx := context.TODO()
			storage := memory.New[*models.Metrics]()
			for k, v := range test.fields {
				storage.Insert(ctx, k, &models.Metrics{
					ID:    k,
					MType: metrics.Counter,
					Delta: &v,
				})
			}
			w := httptest.NewRecorder()
			h := New(&storage, test.f, test.i, metricschecker.New())
			h.Handler()(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			body, err := io.ReadAll(result.Body)
			assert.Equal(t, []byte(test.want.body), body)
			require.NoError(t, err)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
