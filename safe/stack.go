package safe

import (
	"bytes"
	"fmt"
	"runtime"
)

// Stack returns a nicely formatted stack frame, skipping frames.
// Output example:
// runtime error: invalid memory address or nil pointer dereference
//        runtime.gopanic: /usr/local/go/src/runtime/panic.go:965
//        runtime.panicmem: /usr/local/go/src/runtime/panic.go:212
//        runtime.sigpanic: /usr/local/go/src/runtime/signal_unix.go:734
//        github.com/shipengqi/example.v1/cli/internal/action.(*Configuration).Init: /root/gowork/src/cli/internal/action/settings.go:83
//        main.main.func1: /root/gowork/src/cli/main.go:27
//        github.com/spf13/cobra.(*Command).preRun: /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:882
//        github.com/spf13/cobra.(*Command).execute: /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:818
//        github.com/spf13/cobra.(*Command).ExecuteC: /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:960
//        : /root/gowork/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:897
//        main.main: /root/gowork/src/cli/main.go:37
//        runtime.main: /usr/local/go/src/runtime/proc.go:225
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	st := make([]uintptr, 32)
	count := runtime.Callers(skip, st)
	callers := st[:count]
	frames := runtime.CallersFrames(callers)
	for {
		frame, ok := frames.Next()
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s:%d\n", frame.Func.Name(), frame.File, frame.Line)
	}
	return buf.Bytes()
}
