// Package recovery provides a recover function
package recovery

import (
	"github.com/shipengqi/golib/errors"
)

type RecoverFunction func(err error)

// Recovery returns a recover function with error stack added
func Recovery(f RecoverFunction) func() {
	return func() {
		if re := recover(); re != nil {
			var err error
			switch x := re.(type) {
			case string:
				err = errors.New(x)
			case errors.Callers: // avoid duplicate stacks
				err = x
			case error:
				err = errors.New(x.Error())
			default:
				err = errors.New("")
			}
			f(err)
		}
	}
}
