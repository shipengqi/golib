package gssh

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

var ErrNilSession = errors.New("could not start with nil session, use SetSession() to set a session")

type LoggingFunc func(line []byte) error

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments.
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	session *ssh.Session
	ctx     context.Context
}

// newCommand returns the Cmd struct to execute the named program with.
func newCommand(session *ssh.Session, name string, args ...string) *Cmd {
	return &Cmd{
		Path:    name,
		Args:    args,
		session: session,
	}
}

// newCommandContext is like newCommand but includes a context.
func newCommandContext(ctx context.Context, session *ssh.Session, name string, args ...string) *Cmd {
	if ctx == nil {
		panic("nil Context")
	}
	cmd := newCommand(session, name, args...)
	cmd.ctx = ctx
	return cmd
}

// CombinedOutput runs cmd on the remote host and returns its combined
// standard output and standard error.
func (c *Cmd) CombinedOutput() ([]byte, error) {
	if c.session == nil {
		return nil, ErrNilSession
	}
	defer func() { _ = c.session.Close() }()
	return c.execute(func() ([]byte, error) {
		return c.session.CombinedOutput(c.String())
	})
}

// Output runs cmd on the remote host and returns its standard output.
func (c *Cmd) Output() ([]byte, error) {
	if c.session == nil {
		return nil, ErrNilSession
	}
	defer func() { _ = c.session.Close() }()
	return c.execute(func() ([]byte, error) {
		return c.session.Output(c.String())
	})
}

// OutputPipe runs cmd on the remote host and prints its standard output by the given LoggingFunc.
func (c *Cmd) OutputPipe(fn LoggingFunc) error {
	if c.session == nil {
		return ErrNilSession
	}
	defer func() { _ = c.session.Close() }()

	stdout, err := c.session.StdoutPipe()
	if err != nil {
		return err
	}
	err = c.session.Start(c.String())
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, _, rerr := reader.ReadLine()
		if rerr != nil || io.EOF == rerr {
			break
		}
		err = fn(line)
		if err != nil {
			return err
		}
	}
	return c.session.Wait()
}

// Run runs cmd on the remote host.
func (c *Cmd) Run() error {
	if c.session == nil {
		return ErrNilSession
	}
	defer func() { _ = c.session.Close() }()
	_, err := c.execute(func() ([]byte, error) {
		return nil, c.session.Run(c.String())
	})
	return err
}

// Start runs the command on the remote host.
func (c *Cmd) Start() error {
	if c.session == nil {
		return ErrNilSession
	}
	defer func() { _ = c.session.Close() }()
	return c.session.Start(c.String())
}

// String returns a human-readable description of c.
func (c *Cmd) String() string {
	b := new(strings.Builder)
	b.WriteString(c.Path)
	for _, a := range c.Args {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}

// Setenv sets session env vars.
// env specifies the environment of the process.
// Each entry is of the form "key=value", and will be ignored if it is not.
func (c *Cmd) Setenv(env []string) (err error) {
	if c.session == nil {
		return ErrNilSession
	}
	var kv []string
	for _, value := range env {
		kv = strings.SplitN(value, "=", 2)
		if len(kv) < 2 {
			continue
		}
		if err = c.session.Setenv(kv[0], kv[1]); err != nil {
			return
		}
	}

	return
}

// SetSession sets ssh.Session of the command.
func (c *Cmd) SetSession(session *ssh.Session) {
	c.session = session
}

type results struct {
	output []byte
	err    error
}

func (c *Cmd) execute(callback func() ([]byte, error)) ([]byte, error) {
	if c.ctx == nil {
		return callback()
	}
	return c.executeWithCtx(callback)
}

func (c *Cmd) executeWithCtx(callback func() ([]byte, error)) ([]byte, error) {
	done := make(chan results)
	go func() {
		output, err := callback()
		done <- results{
			output: output,
			err:    err,
		}
	}()

	select {
	case <-c.ctx.Done():
		_ = c.session.Signal(ssh.SIGINT)
		return nil, c.ctx.Err()
	case result := <-done:
		return result.output, result.err
	}
}
