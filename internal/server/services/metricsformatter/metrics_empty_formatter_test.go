package metricsformatter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

func TestMetricsEmptyFormatterFormat(t *testing.T) {
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
			s := MetricsEmptyFormatter{}
			str, err := s.Format(tt.m)
			assert.Equal(t, tt.want, str)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
