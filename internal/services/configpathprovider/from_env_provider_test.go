package configpathprovider

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromEnvProvider_Fetch(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			name:    "Test with empty env",
			env:     map[string]string{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "CONFIG=",
			env:     map[string]string{"CONFIG": ""},
			want:    "",
			wantErr: false,
		},
		{
			name:    "CONFIG=/path/to/file",
			env:     map[string]string{"CONFIG": "/path/to/file"},
			want:    "/path/to/file",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			got, err := FromEnvProvider{}.Fetch()
			require.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("%v", err))
			assert.Equal(t, tt.want, got)
		})
	}
}
