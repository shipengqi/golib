package fsutil

import "os"

// owner returns the uid and gid of the given path.
func owner(_ string) (uid, gid int, err error) {
	uid = os.Getuid()
	gid = os.Getgid()
	return
}
