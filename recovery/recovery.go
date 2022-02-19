// Package recovery provides a recover function
package recovery

import (
	"github.com/shipengqi/errors"
)

type RecoverFunction func(err error)

// Recovery returns a recover function with error stack added.
func Recovery(f RecoverFunction) func() {
	return func() {
		if re := recover(); re != nil {
			var err error
			switch x := re.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = errors.New(x.Error()) // deduplicate stacks
			default:
				err = errors.New("")
			}
			f(err)
		}
	}
}
