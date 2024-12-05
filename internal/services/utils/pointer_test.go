package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	t.Run("Test for int", func(t *testing.T) {
		v := 123
		p := Pointer(v)
		assert.Equal(t, v, *p)
	})
	t.Run("Test for float", func(t *testing.T) {
		v := 123.12
		p := Pointer(v)
		assert.Equal(t, v, *p)
	})
	t.Run("Test for string", func(t *testing.T) {
		v := "test"
		p := Pointer(v)
		assert.Equal(t, v, *p)
	})
}
