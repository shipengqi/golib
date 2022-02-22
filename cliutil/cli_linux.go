package cliutil

import (
	"context"
)

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

// ShellExecPipe executes the given command with a pipe that will be connected to the command's
// stdout and stderr when the command starts.
func ShellExecPipe(ctx context.Context, fn LoggingFunc, command string) error {
	if command == "" {
		return ErrInvalidCmd
	}
	return ExecPipe(ctx, fn, "/bin/sh", "-c", command)
}
