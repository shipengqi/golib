package fsutil

import (
	"os"
	"testing"

	"github.com/shipengqi/golib/sysutil"
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

func TestOwner(t *testing.T) {
	if isci() {
		fp := "testdata/ownerfile"
		_, err := os.Create(fp)
		assert.NoError(t, err)

		uid, gid, err := Owner(fp)
		u := sysutil.User()
		assert.NotNil(t, u)

		assert.Equal(t, u.Uid, uid)
		assert.Equal(t, u.Gid, gid)
	}
}

func isci() bool {
	if os.Getenv("CI") == "true" {
		return true
	}
	return false
}
