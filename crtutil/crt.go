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
// with PEM type "CERTIFICATE".
// Deprecated: use ReadFileAsX509 instead.
func ParseCertFile(fpath string) (*x509.Certificate, error) {
	return ReadFileAsX509(fpath)
}

// ReadFileAsX509 read x509.Certificate from the given file.
// The data is expected to be PEM Encoded and contain one certificate
// with PEM type "CERTIFICATE".
func ReadFileAsX509(fpath string) (*x509.Certificate, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return ReadBytesAsX509(bs)
}

// ParseCertBytes parses a single x509.Certificate from the given data.
// The data is expected to be PEM Encoded and contain one certificate
// with PEM type "CERTIFICATE".
// Deprecated: use ReadBytesAsX509 instead.
func ParseCertBytes(data []byte) (*x509.Certificate, error) {
	return ReadBytesAsX509(data)
}

// ReadBytesAsX509 read x509.Certificate from the given data.
// The data is expected to be PEM Encoded and contain one certificate
// with PEM type "CERTIFICATE".
func ReadBytesAsX509(data []byte) (*x509.Certificate, error) {
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
// with PEM type "CERTIFICATE".
// Deprecated: use ReadChainFileAsX509 instead.
func ParseCertChainFile(fpath string) ([]*x509.Certificate, error) {
	return ReadChainFileAsX509(fpath)
}

// ReadChainFileAsX509 read the x509.Certificate chain from the given file.
// The data is expected to be PEM Encoded and contain one of more certificates
// with PEM type "CERTIFICATE".
func ReadChainFileAsX509(fpath string) ([]*x509.Certificate, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return ReadChainBytesAsX509(bs)
}

// ParseCertChainBytes parses x509.Certificate chain from the given data.
// The data is expected to be PEM Encoded and contain one of more certificates
// with PEM type "CERTIFICATE".
// Deprecated: use ReadChainBytesAsX509 instead.
func ParseCertChainBytes(data []byte) ([]*x509.Certificate, error) {
	return ReadChainBytesAsX509(data)
}

// ReadChainBytesAsX509 read x509.Certificate chain from the given data.
// The data is expected to be PEM Encoded and contain one of more certificates
// with PEM type "CERTIFICATE".
func ReadChainBytesAsX509(data []byte) ([]*x509.Certificate, error) {
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

	if len(certs) == 0 {
		return nil, ErrNoPEMData
	}

	return certs, nil
}

// CertToPEM converts a x509.Certificate into a PEM block.
// Deprecated: use EncodeX509ToPEM instead.
func CertToPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

// EncodeX509ToPEM converts a x509.Certificate into a PEM block.
func EncodeX509ToPEM(cert *x509.Certificate, headers map[string]string) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
		Headers: headers,
	})
}

// CertChainToPEM converts a slice of x509.Certificate into PEM block, in the order they are passed.
// Deprecated: use EncodeX509ChainToPEM instead.
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
