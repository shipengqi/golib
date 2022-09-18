package crtutil

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyFile(t *testing.T) {
	prik, err := ReadFileAsSigner("testdata/server-rsa.key")
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}

func TestParseKeyFileWithPass(t *testing.T) {
	prik, err := ReadFileAsSignerWithPass("testdata/server-rsa.key", "")
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}

func TestParseKeyBytes(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/server-rsa-base64.key")
	assert.NoError(t, err)

	prik, err := ReadBytesAsSigner(f, true)
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}
