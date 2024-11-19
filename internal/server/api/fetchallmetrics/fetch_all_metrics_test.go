package fetchallmetrics

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type testStorage struct {
	l map[string]*models.Metrics
}

func (t *testStorage) All(_ context.Context) map[string]*models.Metrics {
	return t.l
}

func floatPointer(val float64) *float64 { return &val }
func intPointer(val int64) *int64       { return &val }

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		s    storage[*models.Metrics]
		want *fetchAllMetricsHandler
	}{
		{
			name: "Test with storage",
			s:    &testStorage{},
			want: &fetchAllMetricsHandler{&testStorage{}},
		},
		{
			name: "Test without storage",
			want: &fetchAllMetricsHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(tt.s))
		})
	}
}

func Test_fetchAllMetricsHandler_Handler(t *testing.T) {
	tests := []struct {
		name       string
		s          storage[*models.Metrics]
		wantInBody []string
		wantCode   int
	}{
		{
			name: "StatusOK with ",
			s: &testStorage{l: map[string]*models.Metrics{
				"gauge145":   {MType: metrics.Gauge, Value: floatPointer(99999.1)},
				"counter111": {MType: metrics.Counter, Delta: intPointer(15345)},
			}},
			wantInBody: []string{"gauge145", "99999.1", "counter111", "15345"},
			wantCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &fetchAllMetricsHandler{s: tt.s}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			h.Handler()(w, r)
			result := w.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			body, _ := io.ReadAll(result.Body)
			for _, val := range tt.wantInBody {
				assert.Contains(t, string(body), val)
			}
		})
	}
}

func Test_fetchAllMetricsHandler_fetchMetrics(t *testing.T) {
	tests := []struct {
		name string
		s    storage[*models.Metrics]
		want map[string]string
	}{
		{
			name: "Test with float values",
			s: &testStorage{l: map[string]*models.Metrics{
				"gauge145": {MType: metrics.Gauge, Value: floatPointer(10.1)},
				"gauge111": {MType: metrics.Gauge, Value: floatPointer(9.1)},
			}},
			want: map[string]string{
				"gauge145": "10.1",
				"gauge111": "9.1",
			},
		},
		{
			name: "Test with float values",
			s: &testStorage{l: map[string]*models.Metrics{
				"counter12": {MType: metrics.Counter, Delta: intPointer(12)},
				"counter19": {MType: metrics.Counter, Delta: intPointer(19)},
			}},
			want: map[string]string{
				"counter12": "12",
				"counter19": "19",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &fetchAllMetricsHandler{s: tt.s}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			assert.Equal(t, tt.want, h.fetchMetrics(r))
		})
	}
}
