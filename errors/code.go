package errors

import (
	"errors"
	"fmt"
	"sync"
)

type causer interface {
	Cause() error
}

type coder interface {
	Code() int
}

const (
	ErrCodeOK      = 0
	ErrCodeUnknown = 1
)

// Reserved errors.
var (
	// OK error code 0.
	OK = Register(ErrCodeOK, "OK")

	// ErrUnknown error code 1.
	ErrUnknown = Register(ErrCodeUnknown, "Unknown Error")
)

// _codes registered codes.
var (
	_codes = make(map[int]struct{})
	mux    = &sync.Mutex{}
)

// Registerf registers an error code with the format specifier.
func Registerf(code int, format string, args ...interface{}) error {
	register(code)
	return &withCode{
		cause: fmt.Errorf(format, args...),
		code:  code,
	}
}

// Register registers an error code with message.
func Register(code int, msg string) error {
	register(code)
	return &withCode{
		cause: errors.New(msg),
		code:  code,
	}
}

// IsCode reports whether any error in err's contains the given code.
func IsCode(err error, code int) bool {
	if v, ok := err.(coder); ok {
		if v.Code() == code {
			return true
		}
	}
	if v, ok := err.(causer); ok {
		err = v.Cause()
		return IsCode(err, code)
	}

	return false
}

func register(code int) {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("[%d] already registered", code))
	}

	mux.Lock()
	defer mux.Unlock()

	_codes[code] = struct{}{}
}

func unregister(code int) {
	if _, ok := _codes[code]; ok {
		mux.Lock()
		defer mux.Unlock()

		delete(_codes, code)
	}
}
