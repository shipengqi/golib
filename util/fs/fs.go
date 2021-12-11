package fs

import (
	"io"
	"os"
	"path/filepath"
)

// HomeDir returns the current user's home directory.
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) (err error) {
	sfd, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() { _ = sfd.Close() }()

	dfd, err := os.OpenFile(dst,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer func() { _ = dfd.Close() }()

	if _, err = io.Copy(dfd, sfd); err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

// Copy copies a file or directory from src to dst.
func Copy(src, dst string) error {
	var (
		err   error
		fds   []os.DirEntry
		sinfo os.FileInfo
	)

	if sinfo, err = os.Stat(src); err != nil {
		return err
	}
	// copies a file
	if !sinfo.IsDir() {
		return CopyFile(src, dst)
	}
	// tries to create dst directory
	if err = os.MkdirAll(dst, sinfo.Mode()); err != nil {
		return err
	}
	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		sfp := filepath.Join(src, fd.Name())
		dfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Copy(sfp, dfp); err != nil {
				return err
			}
		} else {
			if err = CopyFile(sfp, dfp); err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanDir removes all children contained in the given path.
func CleanDir(fpath string) error {
	entries, err := os.ReadDir(fpath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(fpath, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
