package fsutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	t.Run("should be a dir", func(t *testing.T) {
		got := IsDir("testdata/src")
		assert.True(t, got)
	})
	t.Run("should be a file", func(t *testing.T) {
		got := IsDir("testdata/a.txt")
		assert.False(t, got)
	})
}

func TestIsExists(t *testing.T) {
	t.Run("should exist", func(t *testing.T) {
		got := IsExists("testdata/src")
		assert.True(t, got)
	})
	t.Run("should not exist", func(t *testing.T) {
		got := IsExists("testdata/d.txt")
		assert.False(t, got)
	})
}

func TestIsSymlink(t *testing.T) {
	got := IsSymlink("testdata/a.txt")
	assert.False(t, got)
}
