package hashgenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    *hashGenerator
		wantErr bool
	}{
		{
			name:    "Test with empty key",
			key:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with wrong size (3)",
			key:     "123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with wrong size (17)",
			key:     "12345678901234567",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test OK with 16 size key",
			key:     "1234567890123456",
			want:    &hashGenerator{key: "1234567890123456"},
			wantErr: false,
		},
		{
			name:    "Test OK with 24 size key",
			key:     "123456789012345678901234",
			want:    &hashGenerator{key: "123456789012345678901234"},
			wantErr: false,
		},
		{
			name:    "Test OK with 32 size key",
			key:     "12345678901234567890123456789012",
			want:    &hashGenerator{key: "12345678901234567890123456789012"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.key)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.want)
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
		want    []byte
		wantErr bool
	}{
		{
			name:    "Test with empty key",
			key:     "",
			data:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with wrong key 123",
			key:     "123",
			data:    nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := hashGenerator{key: tt.key}
			got, err := g.Generate(tt.data)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
