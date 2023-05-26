package xaes

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	content         = []byte("Hello, World!")
	expected16, _   = base64Decode("DcpY59NkqaesIdfKRiBGAQ==")
	expected24, _   = base64Decode("Qa6PunWL78P2ErIwJ4+TnA==")
	expected32, _   = base64Decode("NHdi3PmU0m6ezdlF2Ueafg==")
	expected16iv, _ = base64Decode("DLviuiaX5JGoUMwpIZIzjg==")
	expected32iv, _ = base64Decode("XocH1UHdlX48PDQB9SMciw==")
	key16           = []byte("1234567891234567")
	key24           = []byte("123456789123456789123456")
	key32           = []byte("12345678912345678912345678912345")
	// iv length must be 16
	iv    = []byte("1234567891234567")
	erriv = []byte("12345678912345678")
)

func TestEncrypt(t *testing.T) {
	got, err := Encrypt(content, key16)
	assert.NoError(t, err)
	assert.Equal(t, expected16, got)

	got, err = Encrypt(content, key24)
	assert.NoError(t, err)
	assert.Equal(t, expected24, got)

	got, err = Encrypt(content, key32)
	assert.NoError(t, err)
	assert.Equal(t, expected32, got)

	got, _ = EncryptWithIV(content, key16, iv)
	assert.NoError(t, err)
	assert.Equal(t, expected16iv, got)

	got, _ = EncryptWithIV(content, key32, iv)
	assert.NoError(t, err)
	assert.Equal(t, expected32iv, got)

	t.Run("Encrypt Error", func(t *testing.T) {
		_, err = EncryptWithIV(content, key16, erriv)
		assert.ErrorContains(t, err, "invalid initial vector")
	})
}

func TestDecrypt(t *testing.T) {

	got, err := Decrypt(expected16, key16)
	assert.NoError(t, err)
	assert.Equal(t, content, got)

	got, err = Decrypt(expected24, key24)
	assert.NoError(t, err)
	assert.Equal(t, content, got)

	got, err = Decrypt(expected32, key32)
	assert.NoError(t, err)
	assert.Equal(t, content, got)

	got, err = DecryptWithIV(expected16iv, key16, iv)
	assert.NoError(t, err)
	assert.Equal(t, content, got)

	got, err = DecryptWithIV(expected32iv, key32, iv)
	assert.NoError(t, err)
	assert.Equal(t, content, got)

	t.Run("Decrypt Error", func(t *testing.T) {
		_, err = DecryptWithIV(expected16, key16, erriv)
		assert.ErrorContains(t, err, "invalid initial vector")
	})
}

func base64Decode(text string) ([]byte, error) {
	data := []byte(text)
	var (
		src    = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err = base64.StdEncoding.Decode(src, data)
	)
	if err != nil {
		return nil, err
	}
	return src[:n], err
}

func base64Encode(data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}
