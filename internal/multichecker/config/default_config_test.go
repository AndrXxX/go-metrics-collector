package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/multichecker/vars"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.Contains(t, config.StaticChecks, vars.StaticSAChecks)
	assert.Contains(t, config.StaticChecks, vars.StaticSTChecks)
	assert.Contains(t, config.StaticChecks, vars.StaticQFChecks)
}
