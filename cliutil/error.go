package cliutil

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidCmd represents a invalid command error.
var ErrInvalidCmd = errors.New("invalid command")

// ErrExit error with stdout and stderr.
type ErrExit struct {
	stdout string
	stderr string
	code   int
}

// Error implements error.
func (e *ErrExit) Error() string {
	b := new(strings.Builder)
	b.WriteString("code: ")
	b.WriteString(strconv.Itoa(e.code))
	if len(e.stderr) > 0 {
		b.WriteString(", ")
		b.WriteString(e.stderr)
	}
	if len(e.stdout) > 0 {
		b.WriteString("\n\t")
		b.WriteString(e.stdout)
	}
	return b.String()
}

// Code returns exit code.
func (e *ErrExit) Code() int {
	return e.code
}

// Stdout returns stdout.
func (e *ErrExit) Stdout() string {
	return e.stdout
}

// Stderr returns stderr.
func (e *ErrExit) Stderr() string {
	return e.stderr
}
