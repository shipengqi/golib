package sysutil

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"unsafe"
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

// Pwd returns the path of the current process.
func Pwd() (home string) {
	var proc string
	var err error
	proc, err = os.Executable()
	if err != nil {
		return
	}
	p, _ := filepath.EvalSymlinks(proc)
	p, _ = filepath.Abs(p)
	home = filepath.Dir(p)
	return
}

// HomeDir returns home directory of current user.
func HomeDir() (home string) {
	var ok bool
	// Unix-like system
	if home, ok = os.LookupEnv("HOME"); ok && home != "" {
		return home
	}

	usr, err := user.Current()
	if nil == err {
		home = usr.HomeDir
	}

	return
}

// IsBigEndian reports whether current os byte order is big endian.
func IsBigEndian() bool {
	v := int32(0x01020304)
	vp := unsafe.Pointer(&v)
	vb := (*byte)(vp)
	// get the value of the pointer
	b := *vb
	// LittleEndian: 04 (03 02 01)
	// BigEndian: 01 (02 03 04)
	return b == 0x01
}

// TmpDir Alias for os.TempDir.
func TmpDir() string {
	return os.TempDir()
}
