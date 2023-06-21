package cliutil

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellExec(t *testing.T) {
	t.Run("say hello", func(t *testing.T) {
		stdout, err := ShellExecContext(context.TODO(), "echo hello, world!")
		assert.NoError(t, err)
		assert.Equal(t, "hello, world!", strings.TrimSpace(stdout))
	})
	t.Run("say hello err", func(t *testing.T) {
		stdout, err := ShellExecContext(context.TODO(), "echooo hello, world!")
		assert.Error(t, err)
		assert.Equal(t, "", strings.TrimSpace(stdout))
		assert.Contains(t, err.Error(), "code: 127")
		assert.Contains(t, err.Error(), "not found")
	})
	t.Run("say hello without ctx", func(t *testing.T) {
		stdout, err := ShellExec("echo hello, world!")
		assert.NoError(t, err)
		assert.Equal(t, "hello, world!", strings.TrimSpace(stdout))
	})
	t.Run("ShellExec ErrExit", func(t *testing.T) {
		_, err := ShellExec("exit 1")
		assert.Error(t, err)
		var exite *ErrExit
		errors.As(err, &exite)
		t.Log(exite.Code(), exite.Stdout(), exite.Stderr())
		assert.Equal(t, 1, exite.Code())
	})
	t.Run("ShellExec ErrInvalidCmd", func(t *testing.T) {
		_, err := ShellExec("")
		assert.ErrorIs(t, err, ErrInvalidCmd)
	})
	t.Run("ShellExecContext ErrInvalidCmd", func(t *testing.T) {
		_, err := ShellExecContext(context.TODO(), "")
		assert.ErrorIs(t, err, ErrInvalidCmd)
	})
	t.Run("ShellExecPipe ErrInvalidCmd", func(t *testing.T) {
		err := ShellExecPipe(context.TODO(), DefaultLoggingFunc, "")
		assert.ErrorIs(t, err, ErrInvalidCmd)
	})
}

func TestShellExecPipe(t *testing.T) {
	testcmd := "echo hello, world!;sleep 1;exit 1"
	testcmd2 := "n=1;while [ $n -le 4 ];do echo $n;((n++));done"
	if os.Getenv("CI") == "true" {
		testcmd = "echo hello, world!;exit 1"
		testcmd2 = "echo 1;echo 2;echo 3;echo 4"
		// t.Skip("Skipped")
	}
	t.Run("exec stdout pipe err", func(t *testing.T) {
		var lines []string
		err := ShellExecPipe(context.TODO(), func(line []byte) {
			lines = append(lines, string(line))
		}, testcmd)
		assert.Equal(t, []string{"hello, world!"}, lines)
		assert.Equal(t, "exit status 1",
			strings.TrimSpace(err.Error()))
	})

	t.Run("exec stdout pipe", func(t *testing.T) {
		var lines []string
		err := ShellExecPipe(context.TODO(), func(line []byte) {
			lines = append(lines, string(line))
		}, testcmd2)
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3", "4"}, lines)
	})
}

func TestExecPipe(t *testing.T) {
	testcmdstderr := "echo 1 1>&2;echo 2 1>&2;echo 3 1>&2;echo 4 1>&2"
	t.Run("exec stderr pipe", func(t *testing.T) {
		var lines []string
		err := ExecErrPipe(context.TODO(), func(line []byte) {
			lines = append(lines, string(line))
		}, "/bin/sh", "-c", testcmdstderr)
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3", "4"}, lines)
	})
}
