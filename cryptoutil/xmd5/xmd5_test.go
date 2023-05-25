package xmd5

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	expected := "65a8e27d8879283831b664bd8b7f0ad4"
	got, err := Encrypt("Hello, World!")
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestEncryptFile(t *testing.T) {
	path := "test.txt"
	expected := "65a8e27d8879283831b664bd8b7f0ad4"
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
