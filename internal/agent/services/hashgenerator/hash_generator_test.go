package hashgenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want *hashGenerator
	}{
		{
			name: "Test with empty key",
			key:  "",
			want: &hashGenerator{key: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, New(tt.key), tt.want)
		})
	}
}

func TestGenerateRandom(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Test with size 10",
			size:    10,
			wantErr: false,
		},
		{
			name:    "Test with size 50",
			size:    50,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateRandom(tt.size)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.size, len(got))
		})
	}
}

func TestHashGeneratorGenerate(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		data    []byte
		want    string
		wantErr bool
	}{
		{
			name:    "Test with empty key",
			key:     "",
			data:    []byte(""),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.key).Generate(tt.data)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
