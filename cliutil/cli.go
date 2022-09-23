package cliutil

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

const defaultPlaceholder = "<>"

// LoggingFunc callback function for logging command output
type LoggingFunc func(line []byte) error

// DefaultLoggingFunc do nothing
func DefaultLoggingFunc(line []byte) error { return nil }

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
// stdout and stderr when the command starts.
func ExecPipe(ctx context.Context, fn LoggingFunc, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() {
		_ = stdout.Close()
	}()
	if err = cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	err = readBuf(reader, fn)
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func readBuf(r *bufio.Reader, fn LoggingFunc) error {
	for {
		if line, _, err := r.ReadLine(); err == nil {
			err = fn(line)
			if err != nil {
				return err
			}
		} else if errors.Is(err, io.EOF) {
			break
		} else {
			return err
		}
	}
	return nil
}
