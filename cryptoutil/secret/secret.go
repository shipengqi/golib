package secret

import "golang.org/x/crypto/bcrypt"

// Encrypt encrypts the plain text with bcrypt.
func Encrypt(v string) (string, error) {
	return EncryptBytes([]byte(v))
}

// EncryptBytes encrypts the bytes with bcrypt.
func EncryptBytes(v []byte) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(v, bcrypt.DefaultCost)
	return string(hashed), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func Compare(encrypted, plain string) error {
	return CompareBytes([]byte(encrypted), []byte(plain))
}

// CompareBytes compares the encrypted bytes with the plain bytes if it's the same.
func CompareBytes(encrypted, plain []byte) error {
	return bcrypt.CompareHashAndPassword(encrypted, plain)
}
