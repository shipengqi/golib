package ssh

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// DefaultTimeout is the default timeout of ssh client connection.
var (
	DefaultUsername = "root"
	DefaultTimeout  = 20 * time.Second
	DefaultPort     = 22
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
	opts *Options
	auth ssh.AuthMethod
}

// New starts a new ssh connection, the host public key must be in known hosts.
func New(opts *Options) (c *Client, err error) {

	// callback, err := DefaultKnownHosts()
	//
	// if err != nil {
	// 	return
	// }
	var auth ssh.AuthMethod
	if opts.UseAgent {
		auth, err = AgentAuth()
		return
	} else {
		auth, err = Auth(opts)
		if err != nil {
			return
		}
	}
	c = &Client{
		opts: opts,
		auth: auth,
	}
	return
}

// Dial starts a client connection to SSH server based on Options.
func (c *Client) Dial(proto string, opts *Options) (*ssh.Client, error) {
	return ssh.Dial(proto, net.JoinHostPort(opts.Addr, strconv.Itoa(opts.Port)),
		&ssh.ClientConfig{
			User:    opts.Username,
			Auth:    []ssh.AuthMethod{c.auth},
			Timeout: opts.Timeout,
		},
	)
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
	// defer func() { _ = session.Close() }()
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

// ExecWithContext starts a new SSH session with context and runs the cmd. It returns CombinedOutput and err if any.
func (c *Client) ExecWithContext(ctx context.Context, name string) ([]byte, error) {
	cmd, err := c.CommandContext(ctx, name)
	if err != nil {
		return nil, err
	}

	return cmd.CombinedOutput()
}

// NewSftp returns new sftp client and error if any.
func (c *Client) NewSftp(opts ...sftp.ClientOption) (*sftp.Client, error) {
	return sftp.NewClient(c.Client, opts...)
}

// Close client net connection.
func (c *Client) Close() error {
	return c.Client.Close()
}
