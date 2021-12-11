package fs

import (
	"os"
	"syscall"
)

// Owner returns the uid and gid of the given path.
func Owner(fpath string) (uid, gid int, err error) {
	uid = os.Getuid()
	gid = os.Getgid()
	if !IsExists(fpath) {
		return
	}
	info, err := os.Stat(fpath)
	if err != nil {
		return
	}
	info.Sys()
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	}
	return
}
