package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	require.NotNil(t, config.Host)
	require.NotNil(t, config.LogLevel)
}
