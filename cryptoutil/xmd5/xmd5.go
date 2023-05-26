// Package xmd5 provides useful functions for MD5 encryption/decryption algorithms.
package xmd5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/shipengqi/golib/convutil"
)

// Encrypt encrypts string with MD5 algorithms.
func Encrypt(v string) (string, error) {
	return EncryptBytes(convutil.S2B(v))
}

// EncryptBytes encrypts data with MD5 algorithms.
func EncryptBytes(v []byte) (encrypted string, err error) {
	h := md5.New()
	if _, err = h.Write(v); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// EncryptFile encrypts the given file content with SHA1 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
