package crtutil

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var (
	ErrUnknownKeyType = errors.New("unknown private key type in PKCS#8 wrapping")
)

// ReadAsX509FromFile read x509.Certificate from the given file.
// The data is expected to be PEM Encoded and contain one or more certificates
// with PEM type "CERTIFICATE".
func ReadAsX509FromFile(fpath string) ([]*x509.Certificate, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return ReadAsX509(bs)
}

// ReadAsX509 read x509.Certificate chain from the given data.
// The data is expected to be PEM Encoded and contain one of more certificates
// with PEM type "CERTIFICATE".
func ReadAsX509(data []byte) ([]*x509.Certificate, error) {
	var (
		certs []*x509.Certificate
		cert *x509.Certificate
		block *pem.Block
		err error
	)

	for len(data) > 0 {
		data = bytes.TrimSpace(data)
		block, data = pem.Decode(data)
		// No PEM data is found
		if block == nil {
			break
		}
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	return certs, nil
}

// EncodeX509ToPEM converts a x509.Certificate into a PEM block.
func EncodeX509ToPEM(cert *x509.Certificate, headers map[string]string) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
		Headers: headers,
	})
}

// EncodeX509ChainToPEM converts a slice of x509.Certificate into PEM block, in the order they are passed.
func EncodeX509ChainToPEM(chain []*x509.Certificate, headers map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	for _, cert := range chain {
		if err := pem.Encode(
			&buf,
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
				Headers: headers,
			},
		); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// IsSelfSigned whether the given x509.Certificate is self-signed.
func IsSelfSigned(cert *x509.Certificate) bool {
	return cert.CheckSignatureFrom(cert) == nil
}
