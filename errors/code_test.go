package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	err := Register(2, "SUCCESS")
	defer unregister(2)
	t.Log(err.Error())
}

func TestRegisterf(t *testing.T) {
	err := Registerf(2, "test %s", "SUCCESS")
	defer unregister(2)
	t.Log(err.Error())
}

func TestRegisterPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, err, "[0] already registered")
		} else {
			t.Fatal("no panic")
		}
	}()
	_ = Register(0, "SUCCESS")
}

func TestRegisterfPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, err, "[0] already registered")
		} else {
			t.Fatal("no panic")
		}
	}()
	_ = Registerf(0, "test %s", "SUCCESS")
}

func TestIsCode(t *testing.T) {
	type run struct {
		expected bool
		code     int
		err      error
	}
	runs := []run{
		{true, 0, OK},
		{false, 0, ErrUnknown},
		{true, 1, ErrUnknown},
		{true, 1, WithCode(ErrUnknown, 2)},
		{true, 1, WithCode(New("test1"), 1)},
		{true, 1, WithCode(WithMessage(New("test2"), "msg2"), 1)},
		{true, 1, WithMessage(ErrUnknown, "msg3")},
		{true, 1, WithMessage(WithCode(ErrUnknown, 2), "msg4")},
		{true, 1, Wrap(ErrUnknown, "msg5")},
		{true, 1, Wrap(WithCode(WithCode(ErrUnknown, 2), 3), "msg6")},
		{true, 2, Wrap(WithCode(WithCode(ErrUnknown, 2), 3), "msg7")},
		{true, 3, Wrap(WithCode(WithCode(ErrUnknown, 2), 3), "msg8")},
	}
	for _, r := range runs {
		got := IsCode(r.err, r.code)
		assert.Equal(t, got, r.expected, fmt.Sprintf("IsCode(%s, %d)", r.err.Error(), r.code))
	}
}
