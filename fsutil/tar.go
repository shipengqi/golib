package fsutil

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Tar create a new archive.
func Tar(src, dst string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = fw.Close() }()

	return tarf(fw, src)
}

// UnTar Extract all files from an archive.
func UnTar(src, dst string) (err error) {
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() { _ = fr.Close() }()

	return untar(fr, dst)
}

// Compress is like Tar but will use gzip to compress.
func Compress(src, dst string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = fw.Close() }()

	gw := gzip.NewWriter(fw)
	defer func() { _ = gw.Close() }()

	return tarf(gw, src)
}

// DeCompress is like UnTar but will use gzip to decompress.
func DeCompress(src, dst string) (err error) {
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() { _ = fr.Close() }()

	// uncompress
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer func() { _ = gr.Close() }()

	return untar(gr, dst)
}

func tarf(writer io.Writer, src string) error {
	tw := tar.NewWriter(writer)
	defer func() { _ = tw.Close() }()
	return filepath.Walk(src, func(filename string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(filename, string(filepath.Separator))
		// write file info
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		// whether info describes a regular file.
		if !info.Mode().IsRegular() {
			return nil
		}
		fr, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() { _ = fr.Close() }()
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
		return nil
	})
}

func untar(reader io.Reader, dst string) error {
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		// Zip Slip Vulnerability
		// See https://cwe.mitre.org/data/definitions/22.html
		if strings.Contains(header.Name, "..") {
			continue
		}
		dstPath := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir: // directory
			if err = MkDirAll(dstPath); err != nil {
				return err
			}
		case tar.TypeReg: // file
			file, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				_ = file.Close()
				return err
			}
			_ = file.Close()
		}
	}
	return nil
}
