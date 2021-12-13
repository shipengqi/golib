package cliutil

import (
	"context"
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
		assert.Equal(t, "code: 127, /bin/sh: echooo: command not found",
			strings.TrimSpace(err.Error()))
	})
	t.Run("say hello without ctx", func(t *testing.T) {
		stdout, err := ShellExec("echo hello, world!")
		assert.NoError(t, err)
		assert.Equal(t, "hello, world!", strings.TrimSpace(stdout))
	})
}

func TestShellExecPipe(t *testing.T) {
	t.Run("exec pipe", func(t *testing.T) {
		var lines []string
		err := ShellExecPipe(context.TODO(), func(line string) {
			// t.Log(line)
			lines = append(lines, line)
		}, "n=1;while [ $n -le 4 ];do echo $n;((n++));done")
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3", "4"}, lines)
	})

	t.Run("exec pipe err", func(t *testing.T) {
		var lines []string
		err := ShellExecPipe(context.TODO(), func(line string) {
			lines = append(lines, line)
		}, "echo hello, world!;sleep 1;exit 1")
		assert.Equal(t, []string{"hello, world!"}, lines)
		assert.Equal(t, "exit status 1",
			strings.TrimSpace(err.Error()))
	})
}
