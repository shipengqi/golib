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
	sftp     *sftp.Client
	opts     *Options
	auth     ssh.AuthMethod
	callback ssh.HostKeyCallback
}

// New creates a Client, the host public key must be in known hosts.
func New(opts *Options) (*Client, error) {
	callback, err := DefaultHostKeyCallback()
	if err != nil {
		return nil, err
	}

	return NewClientWithCallback(opts, callback)
}

// NewInsecure creates a Client that does not verify the server keys.
func NewInsecure(opts *Options) (*Client, error) {
	return NewClientWithCallback(opts, ssh.InsecureIgnoreHostKey())
}

// NewClientWithCallback creates a Client with ssh.HostKeyCallback.
func NewClientWithCallback(opts *Options, callback ssh.HostKeyCallback) (*Client, error) {
	var (
		auth ssh.AuthMethod
		err  error
	)

	auth, err = Auth(opts)
	if err != nil {
		return nil, err
	}

	c := &Client{
		opts:     opts,
		auth:     auth,
		callback: callback,
	}

	err = c.Dial(opts)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Dial starts a client connection to the given SSH server.
func (c *Client) Dial(opts *Options) error {
	cli, err := c.dial(opts)
	if err != nil {
		return err
	}
	c.Client = cli

	return nil
}

func (c *Client) Ping() error {
	_, err := c.dial(c.opts)
	if err != nil {
		return err
	}
	return nil
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

// Upload equivalent to the command `scp <local file> <host>:<remote file>`
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

// Download equivalent to the command `scp <host>:<remote file> <local file>`
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
	return c.Client.Close()
}

func (c *Client) dial(opts *Options) (*ssh.Client, error) {
	return ssh.Dial(DefaultProtocol,
		net.JoinHostPort(opts.Addr, strconv.Itoa(opts.Port)),
		&ssh.ClientConfig{
			User:            opts.Username,
			Auth:            []ssh.AuthMethod{c.auth},
			Timeout:         opts.Timeout,
			HostKeyCallback: c.callback,
		},
	)
}
