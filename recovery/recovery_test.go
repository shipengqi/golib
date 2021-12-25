package recovery

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	pkgerrors "github.com/pkg/errors"
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
		panic(errors.New("pkg panic"))
	})
	t.Run("multi-layers error", func(t *testing.T) {
		defer Recovery(testRecover)()
		panic(pkgerrors.WithStack(errors.New("pkg panic")))
	})
}
