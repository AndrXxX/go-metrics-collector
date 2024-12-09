package fetchallmetrics

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
)

type testStorage struct {
	l map[string]*models.Metrics
}

func (t *testStorage) All(_ context.Context) map[string]*models.Metrics {
	return t.l
}

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
				"gauge145":   {MType: metrics.Gauge, Value: utils.Pointer[float64](99999.1)},
				"counter111": {MType: metrics.Counter, Delta: utils.Pointer[int64](15345)},
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
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			for _, val := range tt.wantInBody {
				assert.Contains(t, string(body), val)
			}
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
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
				"gauge145": {MType: metrics.Gauge, Value: utils.Pointer[float64](10.1)},
				"gauge111": {MType: metrics.Gauge, Value: utils.Pointer[float64](9.1)},
			}},
			want: map[string]string{
				"gauge145": "10.1",
				"gauge111": "9.1",
			},
		},
		{
			name: "Test with float values",
			s: &testStorage{l: map[string]*models.Metrics{
				"counter12": {MType: metrics.Counter, Delta: utils.Pointer[int64](12)},
				"counter19": {MType: metrics.Counter, Delta: utils.Pointer[int64](19)},
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
