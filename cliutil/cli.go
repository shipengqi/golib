package cliutil

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

const defaultPlaceholder = "<>"

// LoggingFunc callback function for logging command output
type LoggingFunc func(line []byte)

// DefaultLoggingFunc do nothing
func DefaultLoggingFunc(line []byte) {}

// RetrieveFlagFromCLI returns value of the given flag from os.Args.
func RetrieveFlagFromCLI(long string, short string) (value string, ok bool) {
	args := os.Args[1:]
	if len(args) < 1 {
		return
	}
	if len(short) == 0 {
		short = defaultPlaceholder // placeholder
	}
	return RetrieveFlag(args, long, short)
}

// RetrieveFlag returns value of the given flag from args.
func RetrieveFlag(args []string, long, short string) (value string, ok bool) {
	var index int
	for k := range args {
		if args[k] == long || args[k] == short {
			index = k + 1
			ok = true
			break
		}
	}
	if index == 0 {
		return
	}
	if len(args) < index+1 {
		return
	}
	if strings.HasPrefix(args[index], "-") {
		return "", ok
	}
	return args[index], ok
}

// ExecContext executes the given command.
func ExecContext(ctx context.Context, command string, args ...string) (output string, err error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
		code   int
	)

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	output = stdout.String()
	code = cmd.ProcessState.ExitCode()

	if err != nil {
		err = &ErrExit{
			stdout: stdout.String(),
			stderr: stderr.String(),
			code:   code,
		}
	}

	return output, err
}

// ExecPipe executes the given command with a pipe that will be connected to the command's
// stdout when the command starts.
func ExecPipe(ctx context.Context, fn LoggingFunc, command string, args ...string) error {
	return pipe(ctx, fn, false, command, args...)
}

// ExecErrPipe executes the given command with a pipe that will be connected to the command's
// stderr when the command starts.
func ExecErrPipe(ctx context.Context, fn LoggingFunc, command string, args ...string) error {
	return pipe(ctx, fn, true, command, args...)
}

func pipe(ctx context.Context, fn LoggingFunc, isstderr bool, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	var (
		rc  io.ReadCloser
		err error
	)
	if isstderr {
		rc, err = cmd.StderrPipe()
	} else {
		rc, err = cmd.StdoutPipe()
	}
	if err != nil {
		return err
	}
	defer func() {
		_ = rc.Close()
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(rc)

	for scanner.Scan() {
		fn(scanner.Bytes())
	}
	return cmd.Wait()
}
