package metricsformatter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

func TestMetricsValueFormatterFormat(t *testing.T) {
	type modelValue struct {
		Value float64
		Delta int64
	}
	tests := []struct {
		name       string
		modelValue *modelValue
		m          *models.Metrics
		want       string
		wantErr    bool
	}{
		{
			name:    "OK Counter Test with empty delta",
			m:       &models.Metrics{ID: "test", MType: metrics.Counter},
			want:    "",
			wantErr: false,
		},
		{
			name:    "OK Gauge Test with empty value",
			m:       &models.Metrics{ID: "test", MType: metrics.Gauge},
			want:    "",
			wantErr: false,
		},
		{
			name:       "OK Counter Test",
			modelValue: &modelValue{Delta: 10},
			m:          &models.Metrics{ID: "test", MType: metrics.Counter},
			want:       "10",
			wantErr:    false,
		},
		{
			name:       "OK Gauge Test",
			modelValue: &modelValue{Value: 10.1},
			m:          &models.Metrics{ID: "test", MType: metrics.Gauge},
			want:       "10.1",
			wantErr:    false,
		},
		{
			name:    "Error unknown metric Test",
			m:       &models.Metrics{ID: "test", MType: "unknown"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MetricsValueFormatter{}
			if tt.modelValue != nil && tt.modelValue.Delta != 0 {
				tt.m.Delta = &tt.modelValue.Delta
			}
			if tt.modelValue != nil && tt.modelValue.Value != 0 {
				tt.m.Value = &tt.modelValue.Value
			}
			str, err := s.Format(tt.m)
			assert.Equal(t, tt.want, str)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
