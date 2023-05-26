package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	tests := []string{
		"123456",
		"xxxxxx",
		"000000",
	}

	for _, v := range tests {
		encrypted, _ := Encrypt(v)
		err := Compare(encrypted, v)
		assert.NoError(t, err)
	}
}
