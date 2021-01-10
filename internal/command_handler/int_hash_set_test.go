package command_handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntHashSet_Hash(t *testing.T) {
	hs := newIntHashSet(3, 2, 1, 9, 10)
	assert.Equal(t, intHashSetHashKey("123910"), hs.Hash())
}

func TestIntHashSet_Contains(t *testing.T) {
	hs := newIntHashSet()
	hs.Store(5)
	hs.Store(4)

	assert.True(t, hs.Contains(5))
	assert.True(t, hs.Contains(4))
	assert.False(t, hs.Contains(120))
}

func TestIntHashSet_Initialise(t *testing.T) {
	hs := newIntHashSet(1, 2, 3)
	assert.True(t, hs.Contains(1))
	assert.True(t, hs.Contains(2))
	assert.True(t, hs.Contains(3))
	assert.False(t, hs.Contains(1234567))
}

func TestIntHashSet_Delete(t *testing.T) {
	hs := newIntHashSet()
	hs.Store(5)

	assert.True(t, hs.Contains(5))

	hs.Delete(5)

	assert.False(t, hs.Contains(5))
}
