package safe

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGo(t *testing.T) {
// 	Go(func() {
// 		t.Log("hello")
// 	})
//
// 	Go(func() {
// 		panic("safe go panic")
// 	})
// }

func TestGoAndWait(t *testing.T) {
	t.Run("should return nil", func(t *testing.T) {
		err := GoAndWait(func() error {
			return nil
		}, func() error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("should return error", func(t *testing.T) {
		err := GoAndWait(func() error {
			return nil
		}, func() error {
			return errors.New("go wait error")
		}, func() error {
			return nil
		})

		assert.Error(t, err, "go wait error")
	})

	t.Run("recover the panic error", func(t *testing.T) {
		err := GoAndWait(func() error {
			return nil
		}, func() error {
			panic("go wait panic")
		}, func() error {
			return nil
		})
		assert.NoError(t, err)
	})
}
