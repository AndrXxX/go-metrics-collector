package metricsidentifier

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestNewURLIdentifier(t *testing.T) {
	tests := []struct {
		name string
		want *urlMetricsIdentifier
	}{
		{
			name: "Test OK NewURLIdentifier",
			want: &urlMetricsIdentifier{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewURLIdentifier())
		})
	}
}

func TestUrlMetricsIdentifierProcess(t *testing.T) {
	type modelValue struct {
		Value float64
		Delta int64
	}
	tests := []struct {
		name       string
		request    string
		vars       map[string]string
		modelValue *modelValue
		want       *models.Metrics
		wantErr    bool
	}{
		{
			name:       "OK Counter test",
			request:    "/update/counter/test/10",
			vars:       map[string]string{vars.MetricType: metrics.Counter, vars.Metric: "test", vars.Value: "10"},
			modelValue: &modelValue{Delta: 10},
			want:       &models.Metrics{ID: "test", MType: metrics.Counter},
			wantErr:    false,
		},
		{
			name:       "OK Gauge test",
			request:    "/update/gauge/test/10.1",
			vars:       map[string]string{vars.MetricType: metrics.Gauge, vars.Metric: "test", vars.Value: "10.1"},
			modelValue: &modelValue{Value: 10.1},
			want:       &models.Metrics{ID: "test", MType: metrics.Gauge},
			wantErr:    false,
		},
		{
			name:    "ОК Gauge test with empty value",
			request: "/update/gauge/test/",
			vars:    map[string]string{vars.MetricType: metrics.Gauge, vars.Metric: "test"},
			want:    &models.Metrics{ID: "test", MType: metrics.Gauge},
			wantErr: false,
		},
		{
			name:    "Error Counter test (bad value)",
			request: "/update/counter/test/aaa",
			vars:    map[string]string{vars.MetricType: metrics.Counter, vars.Metric: "test", vars.Value: "aaa"},
			wantErr: true,
		},
		{
			name:    "Error Gauge test (bad value)",
			request: "/update/gauge/test/aaa",
			vars:    map[string]string{vars.MetricType: metrics.Gauge, vars.Metric: "test", vars.Value: "aaa"},
			wantErr: true,
		},
		{
			name:    "Error test (unknown metric)",
			request: "/update/unknown/test/aaa",
			vars:    map[string]string{vars.MetricType: "unknown", vars.Metric: "test", vars.Value: "aaa"},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := &urlMetricsIdentifier{}
			request := httptest.NewRequest("", test.request, nil)
			ctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			for k, v := range test.vars {
				ctx.URLParams.Add(k, v)
			}
			m, err := i.Process(request)
			require.Equal(t, test.wantErr, err != nil)
			if !test.wantErr {
				assert.Equal(t, test.want.ID, m.ID)
				assert.Equal(t, test.want.MType, m.MType)
				if test.modelValue != nil && test.modelValue.Delta != 0 {
					assert.Equal(t, test.modelValue.Delta, *m.Delta)
				}
				if test.modelValue != nil && test.modelValue.Value != 0 {
					assert.Equal(t, test.modelValue.Value, *m.Value)
				}
			}
		})
	}
}
