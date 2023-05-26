// Package xaes provides useful functions for AES encryption/decryption algorithms, only CBC mode is supported.
package xaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	_IVDefaultValue = []byte("default iv value")
)

// Encrypt encrypts plain text with CBC mode, using the default iv value.
// Note that the key must be 16/24/32 bit length.
func Encrypt(plain []byte, key []byte) ([]byte, error) {
	return EncryptWithIV(plain, key, _IVDefaultValue)
}

// Decrypt decrypts the cipher text in CBC mode, using the default iv value.
// Note that the key must be 16/24/32 bit length.
func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	return DecryptWithIV(cipherText, key, _IVDefaultValue)
}

// EncryptWithIV encrypts plain text with CBC mode.
// Note that the key must be 16/24/32 bit length.
func EncryptWithIV(plain []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// The AES packet length is 128 bits, so blockSize=16
	bs := block.BlockSize()
	if len(iv) != bs {
		return nil, errors.New("invalid initial vector")
	}
	plain = pkcs7Padding(plain, bs)
	// The length of the iv (initial vector) must be equal to the block size
	bm := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plain))
	bm.CryptBlocks(cipherText, plain)

	return cipherText, nil
}

// DecryptWithIV decrypts cipher text with CBC mode.
// Note that the key must be 16/24/32 bit length.
func DecryptWithIV(encrypted []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(iv) != bs {
		return nil, errors.New("invalid initial vector")
	}

	bm := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(encrypted))
	bm.CryptBlocks(origData, encrypted)
	plainText, e := pkcs7UnPadding(origData)
	if e != nil {
		return nil, e
	}
	return plainText, nil
}

func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length <= 0 {
		return nil, errors.New("invalid data len")
	}
	unpadding := int(data[length-1])
	return data[:(length - unpadding)], nil
}
