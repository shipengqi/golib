package cliutil

import (
	"context"
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
	if os.Getenv("CI") == "true" {
		t.Skip("Skipped")
	}
	t.Run("exec pipe", func(t *testing.T) {
		var lines []string
		err := ShellExecPipe(context.TODO(), func(line []byte) error {
			// t.Log(line)
			lines = append(lines, string(line))
			return nil
		}, "n=1;while [ $n -le 4 ];do echo $n;((n++));done")
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3", "4"}, lines)
	})

	t.Run("exec pipe err", func(t *testing.T) {
		var lines []string
		err := ShellExecPipe(context.TODO(), func(line []byte) error {
			lines = append(lines, string(line))
			return nil
		}, "echo hello, world!;sleep 1;exit 1")
		assert.Equal(t, []string{"hello, world!"}, lines)
		assert.Equal(t, "exit status 1",
			strings.TrimSpace(err.Error()))
	})
}
