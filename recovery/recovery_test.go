package recovery

import (
	stderrors "errors"
	"github.com/stretchr/testify/assert"
	"testing"

	pkgerrors "github.com/pkg/errors"

	"github.com/shipengqi/golib/e"
)

func TestRecovery(t *testing.T) {
	var testRecover = func(err error) {
		_, ok := err.(e.Callers)
		assert.True(t, ok)
	}

	t.Run("error string", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic("string panic")
	})
	t.Run("error string", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic("string panic")
	})
	t.Run("std error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(stderrors.New("std panic"))
	})
	t.Run("pkg error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(pkgerrors.New("pkg panic"))
	})
	t.Run("e error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(e.New("pkg panic"))
	})
	t.Run("multi-layers error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(pkgerrors.WithStack(e.New("pkg panic")))
	})
}
