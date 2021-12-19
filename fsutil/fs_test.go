package fsutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	src := "testdata/src"
	dst := "testdata/dst"
	defer func() { _ = os.RemoveAll(dst) }()
	err := Copy(src, dst)
	assert.NoError(t, err)
	entries, err := os.ReadDir(dst)
	assert.NoError(t, err)
	var names []string
	for _, fd := range entries {
		names = append(names, fd.Name())
	}
	assert.ElementsMatch(t, []string{
		"subdir",
		"a.txt",
		"b.txt",
	}, names)
}

func TestCleanDir(t *testing.T) {
	src := "testdata/src"
	dst := "testdata/dst"
	defer func() { _ = os.RemoveAll(dst) }()
	err := Copy(src, dst)
	assert.NoError(t, err)
	entries, err := os.ReadDir(dst)
	assert.NoError(t, err)
	var names []string
	for _, fd := range entries {
		names = append(names, fd.Name())
	}
	assert.ElementsMatch(t, []string{
		"subdir",
		"a.txt",
		"b.txt",
	}, names)
	err = CleanDir(dst)
	assert.NoError(t, err)
	entries, err = os.ReadDir(dst)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(entries))
}

func TestWalk(t *testing.T) {
	src := "testdata/src/a.txt"
	var lines []string
	err := Walk(src, func(line []byte) error {
		lines = append(lines, string(line))
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "22", "333", "4444"}, lines)
}
