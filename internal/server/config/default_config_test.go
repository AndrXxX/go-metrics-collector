package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	require.NotNil(t, config.Host)
}
