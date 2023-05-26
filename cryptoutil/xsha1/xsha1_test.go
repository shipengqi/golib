package xsha1

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	expected := "0a0a9f2a6772942557ab5355d76af442f8f65e01"
	got := Encrypt("Hello, World!")
	assert.Equal(t, expected, got)
}

func TestEncryptFile(t *testing.T) {
	path := "test.txt"
	expected := "0a0a9f2a6772942557ab5355d76af442f8f65e01"
	file, err := os.Create(path)
	assert.NoError(t, err)
	defer func() { _ = os.Remove(path) }()
	defer func() { _ = file.Close() }()

	_, err = file.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	got, err := EncryptFile(path)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
