package requestsender

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

func TestWithXRealIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		params  dto.ParamsDto
		want    dto.ParamsDto
		wantErr bool
	}{
		{
			name:   "Test with empty ip",
			params: dto.ParamsDto{Headers: map[string]string{}},
			want:   dto.ParamsDto{Headers: map[string]string{"X-Real-IP": ""}},
		},
		{
			name:   "Test with ip 192.168.0.1",
			ip:     "192.168.0.1",
			params: dto.ParamsDto{Headers: map[string]string{}},
			want:   dto.ParamsDto{Headers: map[string]string{"X-Real-IP": "192.168.0.1"}},
		},
	}
	for _, tt := range tests {
		f := WithXRealIP(tt.ip)
		require.NoError(t, f(&tt.params))
		require.Equal(t, tt.want, tt.params)
	}
}
