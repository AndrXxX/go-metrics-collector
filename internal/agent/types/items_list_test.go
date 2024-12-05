package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemsList_Add(t *testing.T) {
	t.Run("Test with int", func(t *testing.T) {
		list := ItemsList[int]{15}
		list.Add(10)
		assert.Equal(t, list, ItemsList[int]{15, 10})
	})
	t.Run("Test with string", func(t *testing.T) {
		list := ItemsList[string]{"test1"}
		list.Add("test2")
		assert.Equal(t, list, ItemsList[string]{"test1", "test2"})
	})
}
