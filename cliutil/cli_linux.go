package cliutil

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
)

// LoggingFunc callback function for logging command output
type LoggingFunc func(line []byte) error

// DefaultLoggingFunc do nothing
func DefaultLoggingFunc(line []byte) error { return nil }

// ShellExec executes the given command by shell, errors.g. "ls -al"
func ShellExec(command string) (output string, err error) {
	if command == "" {
		return "", ErrInvalidCmd
	}
	return ExecContext(context.TODO(), "/bin/sh", "-c", command)
}

// ShellExecContext executes the given command by shell.
func ShellExecContext(ctx context.Context, command string) (output string, err error) {
	if command == "" {
		return "", ErrInvalidCmd
	}
	return ExecContext(ctx, "/bin/sh", "-c", command)
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

// ShellExecPipe executes the given command with a pipe that will be connected to the command's
// stdout and stderr when the command starts.
func ShellExecPipe(ctx context.Context, fn LoggingFunc, command string) error {
	if command == "" {
		return ErrInvalidCmd
	}
	return ExecPipe(ctx, fn, "/bin/sh", "-c", command)
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
	err = readPipe(stdout, fn)
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func readPipe(stdout io.ReadCloser, fn LoggingFunc) error {
	reader := bufio.NewReader(stdout)
	return readBuf(reader, fn)
}

func readBuf(r *bufio.Reader, fn LoggingFunc) error {
	for {
		if line, _, err := r.ReadLine(); err == nil {
			err = fn(line)
			if err != nil {
				return err
			}
		} else if err == io.EOF {
			break
		} else {
			return err
		}
	}
	return nil
}
