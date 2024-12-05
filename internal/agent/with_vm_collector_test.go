package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

func TestWithVMCollector(t *testing.T) {
	t.Run("WithVMCollector", func(t *testing.T) {
		opt := WithVMCollector()
		a := agent{c: &config.Config{}}
		assert.Equal(t, 0, len(a.collectors))
		opt(&a)
		assert.Equal(t, 1, len(a.collectors))
	})
}
