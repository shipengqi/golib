package nsenter

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

const (
	SelfPid = "self"
)

// SYS_SETNS syscall allows changing the namespace of the current process.
var SYS_SETNS = map[string]uintptr{
	"386":      346,
	"amd64":    308,
	"arm64":    268,
	"arm":      375,
	"mips":     4344,
	"mipsle":   4344,
	"mips64le": 4344,
	"ppc64":    350,
	"ppc64le":  350,
	"riscv64":  268,
	"s390x":    339,
}[runtime.GOARCH]

var nslist = []string{"ipc", "uts", "net", "pid", "mnt"}

type RecoverFunc func() error

func Set(pid string) (RecoverFunc, error) {
	fds, err := nsfds(pid)
	if err != nil {
		return nil, err
	}
	selfds, err := nsfds(SelfPid)
	if err != nil {
		return nil, err
	}
	runtime.LockOSThread()

	for k := range fds {
		err = SetNs(fds[k].Fd(), 0)
		if err != nil {
			return nil, err
		}
	}

	return func() error {
		for k := range selfds {
			_ = SetNs(selfds[k].Fd(), 0)
			_ = selfds[k].Close()
		}
		runtime.UnlockOSThread()
		for k := range fds {
			_ = fds[k].Close()
		}
		return nil
	}, nil
}

// SetNs sets namespace using syscall.
func SetNs(ns uintptr, nstype int) (err error) {
	_, _, e1 := syscall.Syscall(SYS_SETNS, ns, uintptr(nstype), 0)
	if e1 != 0 {
		err = e1
	}
	return
}

func nsfds(pid string) ([]*os.File, error) {
	var ptrs []*os.File
	for k := range nslist {
		f, err := nsfd(pid, nslist[k])
		if err != nil {
			return nil, err
		}
		ptrs = append(ptrs, f)
	}
	return ptrs, nil
}

func nsfd(pid, ns string) (*os.File, error) {
	var f *os.File
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/%s", pid, ns), os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return f, nil
}
