package hashgenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			assert.Equal(t, Factory(), tt.want)
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
			assert.Equal(t, f.SHA256(), tt.want)
		})
	}
}
