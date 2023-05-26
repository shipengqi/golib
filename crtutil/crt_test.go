package crtutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAsX509FromFile(t *testing.T) {
	tests := []struct {
		title     string
		input     string
		expected  int
		shoulderr bool
	}{
		{
			"read server certificate",
			"testdata/server.crt",
			1,
			false,
		},
		{
			"returns error while reading a non-existent file",
			"testdata/server-non-existent.crt",
			0,
			true,
		},
		{
			"returns 2 CA certificates",
			"testdata/server-ca.crt",
			2,
			false,
		},
		{
			"returns 3 certificates",
			"testdata/server-3layers.crt",
			3,
			false,
		},
		{
			"returns 3 certificates and ignore redundant characters",
			"testdata/server-3layers-withcharacters.crt",
			3,
			false,
		},
	}
	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			certs, err := ReadAsX509FromFile(v.input)
			if v.shoulderr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v.expected, len(certs))
			}
		})
	}

}

func TestParseCertBytes(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected int
	}{
		{
			"empty data, should got 0 certificates",
			[]byte{},
			0,
		},
		{
			"string data, should got 0 certificates",
			[]byte("sdfklhjasdfkjhasdfkjlhas"),
			0,
		},
	}
	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			certs, err := ReadAsX509(v.input)
			assert.NoError(t, err)
			assert.Equal(t, v.expected, len(certs))
		})
	}
}

func TestEncodeX509ChainToPEM(t *testing.T) {
	crts, err := ReadAsX509FromFile("testdata/server-3layers-withcharacters.crt")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(crts))
	got, err := EncodeX509ChainToPEM(crts, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestEncodeX509ToPEM(t *testing.T) {
	crt, err := ReadAsX509FromFile("testdata/server.crt")
	assert.NoError(t, err)
	got := EncodeX509ToPEM(crt[0], nil)
	assert.NotEmpty(t, got)
}

func TestIsSelfSigned(t *testing.T) {
	tests := []struct {
		title    string
		input    string
		expected bool
	}{
		{
			"should not be self-signed",
			"testdata/server.crt",
			false,
		},
		{
			"ca cart should be self-signed",
			"testdata/self-signed.crt",
			true,
		},
		{
			"server cert should be self-signed",
			"testdata/self-signed-not-ca.crt",
			true,
		},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			parsed, err := ReadAsX509FromFile(v.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, parsed)
			got := IsSelfSigned(parsed[0])

			assert.Equal(t, v.expected, got)
		})
	}
}
