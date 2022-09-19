package crtutil

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyFile(t *testing.T) {
	prik, err := ReadAsSignerFromFile("testdata/server-rsa.key")
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}

func TestParseKeyFileWithPass(t *testing.T) {
	prik, err := ReadAsSignerWithPassFromFile("testdata/server-rsa.key", "")
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}

func TestParseKeyBytes(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/server-rsa-base64.key")
	assert.NoError(t, err)

	prik, err := ReadAsSigner(f, true)
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}
