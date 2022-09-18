package crtutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
)

// ParseKeyFile parses an unencrypted crypto.PrivateKey from the given file.
// Deprecated: use ReadFileAsSigner instead.
func ParseKeyFile(fpath string) (crypto.PrivateKey, error) {
	return ReadFileAsSigner(fpath)
}

// ReadFileAsSigner read a crypto.PrivateKey from the given file.
func ReadFileAsSigner(fpath string) (crypto.PrivateKey, error) {
	f, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return ReadBytesAsSigner(f, false)
}

// ParseKeyFileWithPass read a crypto.PrivateKey from the given file.
// Deprecated: use ReadFileAsSignerWithPass instead.
func ParseKeyFileWithPass(keyPath, keyPass string) (crypto.PrivateKey, error) {
	return ReadFileAsSignerWithPass(keyPath, keyPass)
}

// ReadFileAsSignerWithPass read a crypto.PrivateKey from the given file.
func ReadFileAsSignerWithPass(keyPath, keyPass string) (crypto.PrivateKey, error) {
	f, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return readBytesAsSigner(f, []byte(keyPass), false)
}

// ParseKeyBytes read a crypto.PrivateKey from the given data.
// Deprecated: use ReadBytesAsSigner instead.
func ParseKeyBytes(data []byte, isBase64 bool) (crypto.PrivateKey, error) {
	return readBytesAsSigner(data, nil, isBase64)
}

// ReadBytesAsSigner read a crypto.PrivateKey from the given data.
func ReadBytesAsSigner(data []byte, isBase64 bool) (crypto.PrivateKey, error) {
	return readBytesAsSigner(data, nil, isBase64)
}

func readBytesAsSigner(key, keypass []byte, isBase64 bool) (crypto.PrivateKey, error) {
	var err error
	dkeystr := key

	if isBase64 {
		dkeystr, err = base64.StdEncoding.DecodeString(string(key))
		if err != nil {
			return nil, err
		}
	}
	bl, _ := pem.Decode(dkeystr)
	var keyBytes []byte
	if x509.IsEncryptedPEMBlock(bl) && len(keypass) > 0 {
		keyBytes, err = x509.DecryptPEMBlock(bl, keypass)
		if err != nil {
			return nil, err
		}
	} else {
		keyBytes = bl.Bytes
	}

	var pkcs1 *rsa.PrivateKey
	if pkcs1, err = x509.ParsePKCS1PrivateKey(keyBytes); err == nil {
		return pkcs1, nil
	}

	var pkcs8 interface{}
	if pkcs8, err = x509.ParsePKCS8PrivateKey(keyBytes); err == nil {
		switch pkcs8k := pkcs8.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return pkcs8k, nil
		default:
			return nil, ErrUnknownKeyType
		}
	}

	var eck *ecdsa.PrivateKey
	if eck, err = x509.ParseECPrivateKey(keyBytes); err == nil {
		return eck, nil
	}
	return nil, err
}
