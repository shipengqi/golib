package ssh

import (
	"context"
	"net"
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

// New creates a Client, the host public key must be in known hosts.
func New(opts *Options) (*Client, error) {
	var (
		callback ssh.HostKeyCallback
		auth     ssh.AuthMethod
		err      error
	)

	callback, err = DefaultHostKeyCallback()
	if err != nil {
		return nil, err
	}

	auth, err = Auth(opts)
	if err != nil {
		return nil, err
	}

	return &Client{
		opts:     opts,
		auth:     auth,
		callback: callback,
	}, nil
}

// NewInsecure creates a Client that does not verify the server keys.
func NewInsecure(opts *Options) (*Client, error) {
	var (
		auth ssh.AuthMethod
		err  error
	)

	auth, err = Auth(opts)
	if err != nil {
		return nil, err
	}

	return &Client{
		opts:     opts,
		auth:     auth,
		callback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// Dial starts a client connection to the given SSH server.
func (c *Client) Dial(opts *Options) error {
	cli, err := ssh.Dial(DefaultProtocol,
		net.JoinHostPort(opts.Addr, strconv.Itoa(opts.Port)),
		&ssh.ClientConfig{
			User:            opts.Username,
			Auth:            []ssh.AuthMethod{c.auth},
			Timeout:         opts.Timeout,
			HostKeyCallback: c.callback,
		},
	)
	if err != nil {
		return err
	}
	c.Client = cli

	return nil
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

// Close client net connection.
func (c *Client) Close() error {
	return c.Client.Close()
}
