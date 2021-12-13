package crtutil

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var (
	ErrNoPEMData      = errors.New("no pem data is found")
	ErrUnknownKeyType = errors.New("unknown private key type in PKCS#8 wrapping")
)

// ParseCertFile parses x509.Certificate from the given file.
// The data is expected to be PEM Encoded and contain one certificate
// with PEM type "CERTIFICATE"
func ParseCertFile(fpath string) (*x509.Certificate, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return ParseCertBytes(bs)
}

// ParseCertBytes parses a single x509.Certificate from the given data.
// The data is expected to be PEM Encoded and contain one certificate
// with PEM type "CERTIFICATE"
func ParseCertBytes(data []byte) (*x509.Certificate, error) {
	if len(data) == 0 {
		return nil, nil
	}
	bl, _ := pem.Decode(data)
	if bl == nil {
		return nil, ErrNoPEMData
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// ParseCertChainFile parses the x509.Certificate chain from the given file.
// The data is expected to be PEM Encoded and contain one of more certificates
// with PEM type "CERTIFICATE"
func ParseCertChainFile(fpath string) ([]*x509.Certificate, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return ParseCertChainBytes(bs)
}

// ParseCertChainBytes parses x509.Certificate chain from the given data.
// The data is expected to be PEM Encoded and contain one of more certificates
// with PEM type "CERTIFICATE"
func ParseCertChainBytes(data []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	var cert *x509.Certificate
	var block *pem.Block
	var rest []byte
	var err error

	block, rest = pem.Decode(data)
	if block == nil {
		return nil, ErrNoPEMData
	}
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	certs = append(certs, cert)
	for {
		rest = bytes.TrimSpace(rest)
		// This loop terminates because there is no more content
		if len(rest) == 0 {
			break
		}
		block, rest = pem.Decode(rest)
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

// CertToPEM returns a PEM encoded x509 Certificate
func CertToPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

// CertChainToPEM returns a PEM encoded chain of x509 Certificates, in the order they are passed
func CertChainToPEM(chain []*x509.Certificate) ([]byte, error) {
	var buf bytes.Buffer
	for _, cert := range chain {
		if err := pem.Encode(
			&buf,
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			},
		); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
