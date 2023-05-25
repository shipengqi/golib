package xsha256

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"github.com/shipengqi/golib/convutil"
)

// Encrypt encrypts string with SHA1 algorithms.
func Encrypt(v string) string {
	return EncryptBytes(convutil.S2B(v))
}

// EncryptBytes encrypts []byte with SHA1 algorithms.
func EncryptBytes(v []byte) string {
	r := sha1.Sum(v)
	return hex.EncodeToString(r[:])
}

// EncryptFile encrypts the given file content with SHA1 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
