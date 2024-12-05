package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

func TestWithRuntimeCollector(t *testing.T) {
	t.Run("WithRuntimeCollector", func(t *testing.T) {
		opt := WithRuntimeCollector()
		a := agent{c: &config.Config{}}
		assert.Equal(t, 0, len(a.collectors))
		opt(&a)
		assert.Equal(t, 1, len(a.collectors))
	})
}
