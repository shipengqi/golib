package fs

import (
	"os"
)

// IsExists reports whether the given path exists.
func IsExists(fpath string) bool {
	if fpath == "" {
		return false
	}
	_, err := os.Stat(fpath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// IsDir reports whether the given path is a directory.
func IsDir(fpath string) bool {
	if fpath == "" {
		return false
	}
	if fi, err := os.Stat(fpath); err == nil {
		return fi.IsDir()
	}
	return false
}

// IsSymlink reports whether the given path is a symbolic link.
func IsSymlink(fpath string) bool {
	if fpath == "" {
		return false
	}
	if fi, err := os.Stat(fpath); err == nil {
		return fi.Mode()&os.ModeSymlink != 0
	}
	return false
}
