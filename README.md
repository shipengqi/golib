# golib

Common libraries for Go.

[![test](https://github.com/shipengqi/golib/actions/workflows/test.yaml/badge.svg)](https://github.com/shipengqi/golib/actions/workflows/test.yaml)
[![Codecov](https://codecov.io/gh/shipengqi/golib/branch/main/graph/badge.svg?token=SMU4SI304O)](https://codecov.io/gh/shipengqi/golib)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipengqi/golib)](https://goreportcard.com/report/github.com/shipengqi/golib)
[![Release](https://img.shields.io/github/release/shipengqi/golib.svg)](https://github.com/shipengqi/golib/releases)
[![License](https://img.shields.io/github/license/shipengqi/golib)](https://github.com/shipengqi/golib/blob/main/LICENSE)

## Getting Started

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/shipengqi/golib/cliutil"
	"github.com/shipengqi/golib/convutil"
	"github.com/shipengqi/golib/crtutil"
	"github.com/shipengqi/golib/crtutil/tmpl"
	"github.com/shipengqi/golib/cryptoutil/xsha256"
	"github.com/shipengqi/golib/fsutil"
	"github.com/shipengqi/golib/netutil"
	"github.com/shipengqi/golib/retry"
	"github.com/shipengqi/golib/strutil"
	"github.com/shipengqi/golib/sysutil"
)

func main() {
	// --------------------------------------
	// cliutil Examples

	// retrieve value of the given flag from args.
	cliutil.RetrieveFlag(os.Args, "--name", "-n")

	// execute the given command.
	output, _ := cliutil.ExecContext(context.TODO(), "/bin/sh", "-c", "ls -l")
	fmt.Println(output)

	// execute the given command with a pipe.
	pipecmd := "echo 1;echo 2;echo 3;echo 4"
	var lines []string
	_ = cliutil.ExecPipe(context.TODO(), func(line []byte) {
		lines = append(lines, string(line))
	}, "/bin/sh", "-c", pipecmd)
	fmt.Println(lines)
	// output like the following:
	// [1, 2, 3, 4]

	// --------------------------------------
	// convutil Examples

	// convert []byte to string.
	output = convutil.B2S([]byte("abc")) // output: "abc"
	// convert string to []byte.
	_ = convutil.S2B("abc")

	// --------------------------------------
	// crtutil Examples

	// read certificate file
	x509Crt, _ := crtutil.ReadAsX509FromFile("server.crt")

	// converts a slice of x509.Certificate into PEM block, in the order they are passed.
	pemData, _ := crtutil.EncodeX509ChainToPEM(x509Crt, nil)

	// read private key file
	pkey, _ := crtutil.ReadAsSignerFromFile("server.key")

	// output certificate content using the default template
	outputb, _ := tmpl.BuildDefaultCertTemplate(x509Crt[0], true)
	fmt.Println(outputb)
	// output like the following:
	// Serial: 5577006791947779410
	// Valid: 2022-09-23 06:09 UTC to 2032-09-30 06:09 UTC
	// Signature: SHA256-RSA (self-signed)
	// BitLength: 4096
	// Subject Key ID: 6D:E9:2B:2B:1D:59:AB:B5:46:8C:7B:93:C3:49:7E:95:B0:20:E5:4C
	// Basic Constraints: CA:true, pathlen:-1

	
	// --------------------------------------
	// cryptoutil Examples

	// encrypts string with SHA256 algorithms.
	output = xsha256.Encrypt("Hello, World!")
	fmt.Println(output) // dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f
	
	
	// --------------------------------------
	// fsutil Examples

	// copies a file or directory from src to dst.
	_ = fsutil.Copy("testdata/src", "testdata/dst")

	// create a new archive.
	_ = fsutil.Tar("testdata/src", "testdata/dst.tar")

	// extract all files from an archive.
	_ = fsutil.UnTar("testdata/dst.tar", "testdata/src")

	// like Tar but will use gzip to compress.
	_ = fsutil.Compress("testdata/src", "testdata/dst.tgz")

	// like UnTar but will use gzip to decompress.
	_ = fsutil.DeCompress("testdata/dst.tgz", "testdata/src")

	// --------------------------------------
	// netutil Examples

	num := netutil.IPString2Uint("16.187.191.122")
	fmt.Println(num) // 280739706

	// --------------------------------------
	// retry Examples

	var count int
	_ = retry.Times(5).WithInterval(time.Second).Do(func() error {
		count++
		return nil
	})
	fmt.Println(count) // 1

	count = 0
	_ = retry.Times(5).WithInterval(time.Second).Do(func() error {
		count++
		return errors.New("test err")
	})
	fmt.Println(count) // 5

	// --------------------------------------
	// strutil Examples

	// check if str1 contains str2 ignoring case sensitivity
	contains := strutil.ContainsIgnoreCase("STR", "str")
	fmt.Println(contains) // true

	// like strings.ContainsAny but does an "only" instead of "any".
	// If all characters in s are found in chars, the function returns true.
	contains = strutil.ContainsOnly("234234", "0123456789")
	fmt.Println(contains) // true

	// --------------------------------------
	// sysutil Examples

	// returns home directory of current user.
	homedir := sysutil.HomeDir()

	// retrieves the value of the environment variable named
	// by the key
	v := sysutil.EnvOr("TEST_ENV_KEY", "default-value")

	// returns the FQDN of current node.
	fqdn, _ := sysutil.FQDN()
}
```

## Documentation

You can find the docs at [go docs](https://pkg.go.dev/github.com/shipengqi/golib).

## 🔋 JetBrains OS licenses

`golib` had been being developed with **IntelliJ IDEA** under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=golib" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo." width="250" align="middle"></a>