package sys

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// User returns the current user.
func User() *user.User {
	u, err := user.Current()
	if err != nil {
		return nil
	}
	return u
}

// FQDN returns the FQDN of current
func FQDN() (string, error) {
	return fqdn()
}

// IsWindows reports whether the current os is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux reports whether the current os is Linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// ProcessHome returns the path of the current process.
func ProcessHome() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
