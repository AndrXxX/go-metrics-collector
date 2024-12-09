package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{
			name:    "Test with incorrect level",
			level:   "incorrect",
			wantErr: true,
		},
		{
			name:    "Test with level debug",
			level:   "debug",
			wantErr: false,
		},
		{
			name:    "Test with level info",
			level:   "info",
			wantErr: false,
		},
		{
			name:    "Test with level warn",
			level:   "warn",
			wantErr: false,
		},
		{
			name:    "Test with level error",
			level:   "error",
			wantErr: false,
		},
		{
			name:    "Test with level fatal",
			level:   "fatal",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, Initialize(tt.level) != nil)
			if !tt.wantErr {
				require.Equal(t, tt.level, Log.Level().String())
			}
		})
	}
}
