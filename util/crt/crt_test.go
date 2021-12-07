package crt

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCrtFile(t *testing.T) {
	crt, err := ParseCertFile("testdata/server.crt")
	if err != nil {
		t.Fatal(err)
	}
	printCrt(t, crt, "server")
}

func TestParseCrtSetFile(t *testing.T) {
	crts, err := ParseCertChainFile("testdata/server-ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(crts))
	for i, crt := range crts {
		printCrt(t, crt, fmt.Sprintf("server-ca-%d", i))
	}

	crts, err = ParseCertChainFile("testdata/server-3layers.crt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(crts))
	for i, crt := range crts {
		printCrt(t, crt, fmt.Sprintf("server-3layers-%d", i))
	}

	crts, err = ParseCertChainFile("testdata/server-3layers-withcharacters.crt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(crts))
	for i, crt := range crts {
		printCrt(t, crt, fmt.Sprintf("server-3layers-withcharacters-%d", i))
	}
}

func TestCertChainToPEM(t *testing.T) {
	crts, err := ParseCertChainFile("testdata/server-3layers-withcharacters.crt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(crts))
	got, err := CertChainToPEM(crts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", got)
}

func TestCertToPEM(t *testing.T) {
	crt, err := ParseCertFile("testdata/server.crt")
	if err != nil {
		t.Fatal(err)
	}
	got := CertToPEM(crt)
	t.Logf("\n%s", got)
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
