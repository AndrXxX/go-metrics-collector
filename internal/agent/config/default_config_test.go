package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	require.NotNil(t, config.Common)
	assert.NotEmpty(t, config.Common.Host)
	require.NotNil(t, config.Intervals)
	assert.NotNil(t, config.Intervals.PollInterval)
	assert.NotNil(t, config.Intervals.ReportInterval)
	assert.NotNil(t, config.Intervals.SleepInterval)
	require.NotNil(t, config.Metrics)
	assert.Contains(t, config.Metrics, metrics.Alloc)
}
