package jsontime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	jsonEncoder := func(v any) []byte {
		val, _ := json.Marshal(v)
		return val
	}
	tests := []struct {
		name    string
		b       []byte
		wantErr bool
		want    time.Duration
	}{
		{
			name:    "Test with error on unmarshal",
			b:       []byte("{{fd}"),
			wantErr: true,
			want:    time.Duration(0),
		},
		{
			name:    "Test with string value",
			b:       jsonEncoder("1s"),
			wantErr: false,
			want:    1 * time.Second,
		},
		{
			name:    "Test with error on decode string value",
			b:       jsonEncoder("1..11s"),
			wantErr: true,
		},
		{
			name:    "Test with float value",
			b:       jsonEncoder(1.1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Duration{Duration: *new(time.Duration)}
			err := d.UnmarshalJSON(tt.b)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
