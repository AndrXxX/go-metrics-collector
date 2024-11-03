package metricsvaluesetter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

func TestGaugeValueSetterSet(t *testing.T) {
	type args struct {
		m     *models.Metrics
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantValue float64
		wantErr   bool
	}{
		{
			name:      "Test OK with 10.1",
			args:      args{m: &models.Metrics{}, value: "10.1"},
			wantValue: 10.1,
			wantErr:   false,
		},
		{
			name:      "Test OK with 10",
			args:      args{m: &models.Metrics{}, value: "10"},
			wantValue: 10,
			wantErr:   false,
		},
		{
			name:    "Test Error with aaa",
			args:    args{m: &models.Metrics{}, value: "aaa"},
			wantErr: true,
		},
		{
			name:    "Test Error with empty string",
			args:    args{m: &models.Metrics{}, value: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setter := &gaugeValueSetter{}
			err := setter.Set(tt.args.m, tt.args.value)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.InDelta(t, tt.wantValue, *tt.args.m.Value, 0.0001)
			}
		})
	}
}
