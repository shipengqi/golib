package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	t.Run("should be a dir", func(t *testing.T) {
		got := IsDir("testdata/src")
		assert.Equal(t, true, got)
	})
	t.Run("should be a file", func(t *testing.T) {
		got := IsDir("testdata/a.txt")
		assert.Equal(t, false, got)
	})
}

func TestIsExists(t *testing.T) {
	t.Run("should exist", func(t *testing.T) {
		got := IsExists("testdata/src")
		assert.Equal(t, true, got)
	})
	t.Run("should not exist", func(t *testing.T) {
		got := IsExists("testdata/d.txt")
		assert.Equal(t, false, got)
	})
}

func TestIsSymlink(t *testing.T) {
	got := IsSymlink("testdata/a.txt")
	assert.Equal(t, false, got)
}
