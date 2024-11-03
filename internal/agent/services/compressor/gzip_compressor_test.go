package compressor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGzipCompressor_Compress(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "Test with empty data",
			data:    nil,
			wantErr: false,
		},
		{
			name:    "Test with data",
			data:    []byte("test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GzipCompressor{}
			_, err := c.Compress(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
