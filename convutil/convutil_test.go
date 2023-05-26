package convutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		s string
		b []byte
	}{
		{"123", []byte("123")},
		{"abc", []byte("abc")},
	}

	for _, v := range tests {
		gotb := S2B(v.s)
		assert.Equal(t, v.b, gotb)
		gots := B2S(gotb)
		assert.Equal(t, v.s, gots)
	}
}
