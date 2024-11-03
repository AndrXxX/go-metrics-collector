package hashgenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		name string
		want *hashGeneratorFactory
	}{
		{
			name: "Test Ok",
			want: &hashGeneratorFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Factory())
		})
	}
}

func TestFactorySHA256(t *testing.T) {
	tests := []struct {
		name string
		want *sha256Generator
	}{
		{
			name: "Test Ok",
			want: &sha256Generator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factory()
			assert.Equal(t, tt.want, f.SHA256())
		})
	}
}
