package xsha256

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	expected := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	got := Encrypt("Hello, World!")
	assert.Equal(t, expected, got)
}

func TestEncryptFile(t *testing.T) {
	path := "test.txt"
	expected := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
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
