package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/vars"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.Contains(t, config.StaticAnalyzers, vars.StaticSAAnalyzers)
	assert.Contains(t, config.StaticAnalyzers, vars.StaticSTAnalyzers)
	assert.Contains(t, config.StaticAnalyzers, vars.StaticQFAnalyzers)
	assert.Contains(t, config.ExcludeStaticAnalyzers, vars.StaticST1000Analyzer)
}
