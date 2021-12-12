package fs

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
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	entries, err = os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(entries))
}
