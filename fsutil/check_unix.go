//go:build linux || darwin

package fsutil

import (
	"os"
	"syscall"
)

func owner(fpath string) (uid, gid int, err error) {
	uid = os.Getuid()
	gid = os.Getgid()
	if !IsExists(fpath) {
		return
	}
	info, err := os.Stat(fpath)
	if err != nil {
		return
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	}
	return
}
