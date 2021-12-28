// Package gssh provides a simple SSH client for Go.
package gssh

import (
	"context"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	// DefaultUsername default user of ssh client connection.
	DefaultUsername = "root"

	// DefaultTimeout default timeout of ssh client connection.
	DefaultTimeout = 20 * time.Second

	// DefaultPort default port of ssh client connection.
	DefaultPort     = 22
	DefaultProtocol = "tcp"
)

// Options for SSH Client.
type Options struct {
	Username   string
	Password   string
	Key        string
	Passphrase string
	Addr       string
	Port       int
	UseAgent   bool
	Timeout    time.Duration
}

// NewOptions creates an Options with default parameters.
func NewOptions() *Options {
	return &Options{
		Username: DefaultUsername,
		Port:     DefaultPort,
		Timeout:  DefaultTimeout,
	}
}

// Client SSH client.
type Client struct {
	*ssh.Client
	opts     *Options
	auth     ssh.AuthMethod
	callback ssh.HostKeyCallback
}

// NewDefault creates a Client with DefaultHostKeyCallback, the host public key must be in known hosts.
func NewDefault(opts *Options) (*Client, error) {
	callback, err := DefaultHostKeyCallback()
	if err != nil {
		return nil, err
	}
	cli, err := New(opts)
	if err != nil {
		return nil, err
	}
	cli.WithHostKeyCallback(callback)
	return cli, nil
}

// NewInsecure creates a Client that does not verify the server keys.
func NewInsecure(opts *Options) (*Client, error) {
	cli, err := New(opts)
	if err != nil {
		return nil, err
	}
	cli.WithHostKeyCallback(ssh.InsecureIgnoreHostKey())
	return cli, nil
}

// New creates a Client without ssh.HostKeyCallback.
func New(opts *Options) (*Client, error) {
	var (
		auth ssh.AuthMethod
		err  error
	)

	auth, err = Auth(opts)
	if err != nil {
		return nil, err
	}

	c := &Client{
		opts: opts,
		auth: auth,
		callback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return c, nil
}

// WithHostKeyCallback sets ssh.HostKeyCallback of Client.
func (c *Client) WithHostKeyCallback(callback ssh.HostKeyCallback) *Client {
	c.callback = callback
	return c
}

// Dial starts a client connection to the given SSH server.
func (c *Client) Dial() error {
	cli, err := c.dial()
	if err != nil {
		return err
	}
	c.Client = cli

	return nil
}

// Ping alias for Dial.
func (c *Client) Ping() error {
	return c.Dial()
}

// CombinedOutput runs cmd on the remote host and returns its combined
// standard output and standard error.
func (c *Client) CombinedOutput(command string) ([]byte, error) {
	session, err := c.NewSession()
	if err != nil {
		return nil, err
	}

	defer func() { _ = session.Close() }()

	return session.CombinedOutput(command)
}

// CombinedOutputContext is like CombinedOutput but includes a context.
//
// The provided context is used to kill the process (by calling
// os.Process.Kill) if the context becomes done before the command
// completes on its own.
func (c *Client) CombinedOutputContext(ctx context.Context, command string) ([]byte, error) {
	cmd, err := c.CommandContext(ctx, command)
	if err != nil {
		return nil, err
	}
	return cmd.CombinedOutput()
}

// Command returns the Cmd struct to execute the named program with
// the given arguments.
//
// It sets only the Path and Args in the returned structure.
func (c *Client) Command(name string, args ...string) (*Cmd, error) {
	session, err := c.NewSession()
	if err != nil {
		return nil, err
	}
	return newCommand(session, name, args...), nil
}

// CommandContext is like Command but includes a context.
//
// The provided context is used to kill the process (by calling
// os.Process.Kill) if the context becomes done before the command
// completes on its own.
func (c *Client) CommandContext(ctx context.Context, name string, args ...string) (*Cmd, error) {
	session, err := c.NewSession()
	if err != nil {
		return nil, err
	}
	return newCommandContext(ctx, session, name, args...), nil
}

// NewSftp returns new sftp client and error if any.
func (c *Client) NewSftp(opts ...sftp.ClientOption) (*sftp.Client, error) {
	return sftp.NewClient(c.Client, opts...)
}

// Upload equivalent to the command `scp <local file> <host>:<remote file>`.
func (c *Client) Upload(lpath, rpath string) (err error) {
	local, err := os.Open(lpath)
	if err != nil {
		return
	}
	defer func() { _ = local.Close() }()

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer func() { _ = ftp.Close() }()

	remote, err := ftp.Create(rpath)
	if err != nil {
		return
	}
	defer func() { _ = remote.Close() }()

	_, err = io.Copy(remote, local)
	return
}

// Download equivalent to the command `scp <host>:<remote file> <local file>`.
func (c *Client) Download(rpath, lpath string) (err error) {
	local, err := os.Create(lpath)
	if err != nil {
		return
	}
	defer func() { _ = local.Close() }()

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer func() { _ = ftp.Close() }()

	remote, err := ftp.Open(rpath)
	if err != nil {
		return
	}
	defer func() { _ = remote.Close() }()

	if _, err = io.Copy(local, remote); err != nil {
		return
	}

	return local.Sync()
}

// Close client net connection.
func (c *Client) Close() error {
	if c.Client != nil {
		return c.Client.Close()
	}
	return nil
}

func (c *Client) dial() (*ssh.Client, error) {
	return ssh.Dial(DefaultProtocol,
		net.JoinHostPort(c.opts.Addr, strconv.Itoa(c.opts.Port)),
		&ssh.ClientConfig{
			User:            c.opts.Username,
			Auth:            []ssh.AuthMethod{c.auth},
			Timeout:         c.opts.Timeout,
			HostKeyCallback: c.callback,
		},
	)
}

func Ping(addr, user, password, key string) error {
	cli, err := NewInsecure(&Options{
		Username: user,
		Password: password,
		Key:      key,
		Addr:     addr,
		Port:     DefaultPort,
	})
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()
	return cli.Ping()
}
