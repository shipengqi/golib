package recovery

import (
	stderrors "errors"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/shipengqi/errors"
)

func TestRecovery(t *testing.T) {
	var testRecover = func(err error) {
		_, ok := err.(errors.Callers)
		assert.True(t, ok)
	}

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
	t.Run("errors error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(errors.New("golib err panic"))
	})
	t.Run("multi-layers error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(pkgerrors.WithStack(errors.New("golib and pkg err panic")))
	})
}
