package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	require.NotNil(t, config.Host)
	require.NotNil(t, config.LogLevel)
}

func ExampleNewConfig() {
	config := NewConfig()

	fmt.Println(config.Host)
	fmt.Println(config.LogLevel)
	fmt.Println(config.StoreInterval)

	// Output:
	// localhost:8080
	// info
	// 300
}
