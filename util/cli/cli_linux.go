package cli

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"sync"
)

// LoggingCallback callback function for logging command output
type LoggingCallback func(line string)

// DefaultLoggingCallback do nothing
func DefaultLoggingCallback(line string) {}

// ShellExec executes the given command by shell, e.g. "ls -al"
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
func ShellExecPipe(ctx context.Context, callback LoggingCallback, command string) error {
	if command == "" {
		return ErrInvalidCmd
	}
	return ExecPipe(ctx, callback, "/bin/sh", "-c", command)
}

// ExecPipe executes the given command with a pipe that will be connected to the command's
// stdout and stderr when the command starts.
func ExecPipe(ctx context.Context, callback LoggingCallback, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	err = readPipe(stdout, stderr, callback)
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func readPipe(stdout, stderr io.ReadCloser, callback LoggingCallback) error {
	oReader := bufio.NewReader(stdout)
	eReader := bufio.NewReader(stderr)
	wg := &sync.WaitGroup{}

	wg.Add(2)
	var oErr, eErr error
	go func() {
		defer wg.Done()
		oErr = readBuf(oReader, callback)
	}()
	go func() {
		defer wg.Done()
		eErr = readBuf(eReader, callback)
	}()
	wg.Wait()
	if oErr != nil {
		return eErr
	}
	if eErr != nil {
		return eErr
	}
	return nil
}

func readBuf(r *bufio.Reader, callback LoggingCallback) error {
	for {
		if line, _, err := r.ReadLine(); err == nil {
			callback(string(line))
		} else if err == io.EOF {
			break
		} else {
			return err
		}
	}
	return nil
}
