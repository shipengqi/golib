package cliutil

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidCmd = errors.New("invalid command")
)

type ErrExit struct {
	stdout string
	stderr string
	code   int
}

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

func (e *ErrExit) Code() int {
	return e.code
}

func (e *ErrExit) Stdout() string {
	return e.stdout
}

func (e *ErrExit) Stderr() string {
	return e.stderr
}
