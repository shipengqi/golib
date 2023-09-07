package fsutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTar(t *testing.T) {
	src := "testdata/src"
	dst := "testdata/dst.tar"
	defer func() { _ = os.RemoveAll(dst) }()
	err := Tar(src, dst)
	assert.NoError(t, err)

	unsrc := dst
	undst := "testdata/src2"
	defer func() { _ = os.RemoveAll(undst) }()
	err = UnTar(unsrc, undst)
	assert.NoError(t, err)

	f, err := os.Stat("testdata/src2/testdata/src/subdir")
	assert.NoError(t, err)
	assert.True(t, f.IsDir())

	f, err = os.Stat("testdata/src2/testdata/src/a.txt")
	assert.NoError(t, err)
	assert.False(t, f.IsDir())

	f, err = os.Stat("testdata/src2/testdata/src/b.txt")
	assert.NoError(t, err)
	assert.False(t, f.IsDir())

	f, err = os.Stat("testdata/src2/testdata/src/subdir/c.txt")
	assert.NoError(t, err)
	assert.False(t, f.IsDir())
}

func TestTarFile(t *testing.T) {
	src := "testdata/src/subdir/c.txt"
	dst := "testdata/dst.tar"
	defer func() { _ = os.RemoveAll(dst) }()
	err := Tar(src, dst)
	assert.NoError(t, err)

	unsrc := dst
	undst := "testdata/src2"
	defer func() { _ = os.RemoveAll(undst) }()
	err = UnTar(unsrc, undst)
	assert.NoError(t, err)
	f, err := os.Stat("testdata/src2/testdata/src/subdir")
	assert.NoError(t, err)
	assert.True(t, f.IsDir())

	f, err = os.Stat("testdata/src2/testdata/src/subdir/c.txt")
	assert.NoError(t, err)
	assert.False(t, f.IsDir())
}

func TestCompress(t *testing.T) {
	src := "testdata/src"
	dst := "testdata/dst.tar.gz"
	defer func() { _ = os.RemoveAll(dst) }()
	err := Compress(src, dst)
	assert.NoError(t, err)

	unsrc := dst
	undst := "testdata/src2"
	defer func() { _ = os.RemoveAll(undst) }()
	err = DeCompress(unsrc, undst)
	assert.NoError(t, err)
}
