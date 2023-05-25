package xaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	_IVDefaultValue = []byte("xaes default iv value")
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
	// key的长度必须是16、24或者32字节，分别用于选择AES-128, AES-192, or AES-256
	// AES 分组长度为 128 位，所以 blockSize=16，单位字节
	bs := block.BlockSize()
	plain = PKCS5Padding(plain, bs)
	bm := cipher.NewCBCEncrypter(block, iv) // 初始向量的长度必须等于块block的长度16字节
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
	if len(encrypted) < bs {
		return nil, errors.New("cipher too short")
	}
	if len(encrypted)%bs != 0 {
		return nil, errors.New("cipher is not a multiple of the block size")
	}

	bm := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(encrypted))
	bm.CryptBlocks(origData, encrypted)
	plainText, e := PKCS5UnPadding(origData, bs)
	if e != nil {
		return nil, e
	}
	return plainText, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, errors.New("invalid block len")
	}

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

	unpadding := int(src[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, errors.New("invalid padding")
	}

	padding := src[length-unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}
