//go:build go1.20

package convutil

import (
	"unsafe"
)

// B2S convert []byte to string.
// See https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// S2B convert string to []byte.
// See https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
