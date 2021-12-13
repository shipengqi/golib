package recovery

import (
	"github.com/shipengqi/golib/e"
	"testing"

	stderrors "errors"

	pkgerrors "github.com/pkg/errors"
)

func TestRecovery(t *testing.T) {
	var testRecover = func(err error) {
		t.Logf("[recover panic]: %+v",
			err)
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
