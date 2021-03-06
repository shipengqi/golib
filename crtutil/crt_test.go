package crtutil

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCrtFile(t *testing.T) {
	_, err := ParseCertFile("testdata/server.crt")
	assert.NoError(t, err)
	// printCrt(t, crt, "server")

	t.Run("parse error with empty data", func(t *testing.T) {
		_, err = ParseCertFile("testdata/server-fail.crt")
		assert.Error(t, err)
	})
}

func TestParseCertBytes(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		_, err := ParseCertBytes([]byte{})
		assert.NoError(t, err)
	})
	t.Run("ErrNoPEMData", func(t *testing.T) {
		_, err := ParseCertBytes([]byte("sdfklhjasdfkjhasdfkjlhas"))
		assert.ErrorIs(t, err, ErrNoPEMData)
	})
}

func TestParseCrtSetFile(t *testing.T) {
	crts, err := ParseCertChainFile("testdata/server-ca.crt")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(crts))

	crts, err = ParseCertChainFile("testdata/server-3layers.crt")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(crts))

	crts, err = ParseCertChainFile("testdata/server-3layers-withcharacters.crt")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(crts))

	t.Run("parse error with empty data", func(t *testing.T) {
		_, err = ParseCertChainFile("testdata/server-fail.crt")
		assert.Error(t, err)
	})
}

func TestParseCertChainBytes(t *testing.T) {
	t.Run("ErrNoPEMData", func(t *testing.T) {
		_, err := ParseCertChainBytes([]byte("sdfklhjasdfkjhasdfkjlhas"))
		assert.ErrorIs(t, err, ErrNoPEMData)
	})
}

func TestCertChainToPEM(t *testing.T) {
	crts, err := ParseCertChainFile("testdata/server-3layers-withcharacters.crt")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(crts))
	got, err := CertChainToPEM(crts)
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestCertToPEM(t *testing.T) {
	crt, err := ParseCertFile("testdata/server.crt")
	assert.NoError(t, err)
	got := CertToPEM(crt)
	assert.NotEmpty(t, got)
}

func printCrt(t *testing.T, cert *x509.Certificate, name string) {
	t.Log("")
	t.Logf("%s Certificate Information:", name)
	t.Logf("  Issuer: %s", cert.Issuer)
	t.Logf("  NotBefore: %s", cert.NotBefore.String())
	t.Logf("  NotAfter: %s", cert.NotAfter.String())
	t.Logf("  Subject: %s", cert.Subject)

	dnsStr := strings.Join(cert.DNSNames, ",")
	ipBuf := new(bytes.Buffer)

	for k := range cert.IPAddresses {
		if k == 0 {
			_, _ = fmt.Fprintf(ipBuf, "%s", cert.IPAddresses[k].String())
		} else {
			_, _ = fmt.Fprintf(ipBuf, ", %s", cert.IPAddresses[k].String())
		}
	}
	t.Logf("  DNSNames: %s", dnsStr)
	t.Logf("  IPAddresses: %s", ipBuf.String())
	t.Logf("  KeyUsage: %v", cert.KeyUsage)
	t.Logf("  ExtKeyUsage: %v", cert.ExtKeyUsage)
	t.Logf("  IsCA: %v", cert.IsCA)
}
