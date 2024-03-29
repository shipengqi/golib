//go:build !go1.20

package convutil

import (
	"reflect"
	"unsafe"
)

// B2S convert []byte to string.
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B convert string to []byte.
func S2B(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len

	return b
}
