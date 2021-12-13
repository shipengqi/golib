package crtutil

import (
	"crypto/rsa"
	"testing"
)

func TestParseKeyFile(t *testing.T) {
	prik, err := ParseKeyFile("testdata/server-rsa.key")
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := prik.(*rsa.PrivateKey); !ok {
		t.Fatal("not rsa key")
	}
}
