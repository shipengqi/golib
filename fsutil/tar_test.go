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
