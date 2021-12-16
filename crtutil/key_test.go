package crtutil

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyFile(t *testing.T) {
	prik, err := ParseKeyFile("testdata/server-rsa.key")
	assert.NoError(t, err)

	_, ok := prik.(*rsa.PrivateKey)
	assert.True(t, ok)
}
