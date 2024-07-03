package metricstringifier

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsEmptyStringifierString(t *testing.T) {
	tests := []struct {
		name    string
		m       *models.Metrics
		want    string
		wantErr bool
	}{
		{
			name:    "OK Test",
			m:       &models.Metrics{},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MetricsEmptyStringifier{}
			str, err := s.String(tt.m)
			assert.Equal(t, tt.want, str)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
