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

// go test -v -run=none -bench=^BenchmarkBytesConv -benchmem=true

var testString = "This string is used to benchmark tests."
var testBytes = []byte(testString)

func rawB2S(b []byte) string {
	return string(b)
}

func rawS2B(s string) []byte {
	return []byte(s)
}

func BenchmarkBytesConvBytesToStrRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawB2S(testBytes)
	}
}

func BenchmarkBytesConvBytesToStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		B2S(testBytes)
	}
}

func BenchmarkBytesConvStrToBytesRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawS2B(testString)
	}
}

func BenchmarkBytesConvStrToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		S2B(testString)
	}
}
