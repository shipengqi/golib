package recovery

import "github.com/shipengqi/golib/e"

type RecoverFunction func(err error)

// Recovery returns a recover function with error stack added
func Recovery(f RecoverFunction) func() {
	return func() {
		if re := recover(); re != nil {
			var err error
			switch x := re.(type) {
			case string:
				err = e.New(x)
			case e.Callers: // avoid duplicate stacks
				err = x
			case error:
				err = e.WithStack(x)
			default:
				err = e.New("")
			}
			f(err)
		}
	}
}
