package gssh

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	addr   string
	user   string
	passwd string
	key    string
)

func TestGSSH(t *testing.T) {
	t.Run("TestPassAuth", secure(t, passAuthTest))
}

func TestGSSHInsecure(t *testing.T) {
	t.Run("TestPassAuth", insecure(t, passAuthTest))
}

func insecure(t *testing.T, callback func(t *testing.T, cli *Client)) func (t *testing.T) {
	opts := NewOptions()
	opts.Username = user
	opts.Password = passwd
	opts.Addr = addr
	opts.Key = key

	cli, err := NewInsecure(opts)
	if err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		callback(t, cli)
	}
}

func secure(t *testing.T, callback func(t *testing.T, cli *Client)) func (t *testing.T) {
	opts := NewOptions()
	opts.Username = user
	opts.Password = passwd
	opts.Addr = addr
	opts.Key = key

	cli, err := NewClientWithCallback(opts, AutoFixedHostKeyCallback)
	if err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		callback(t, cli)
	}
}

func passAuthTest(t *testing.T, cli *Client) {
	cmd, err := cli.Command("echo", "Hello, world!")
	if err != nil {
		t.Fatal(err)
	}

	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Hello, world!\n", string(output))
}

func TestMain(m *testing.M) {
	flag.StringVar(&addr, "addr", "localhost", "The host of ssh")
	flag.StringVar(&user, "user", "root", "The username of client")
	flag.StringVar(&passwd, "pass", "", "The password of user")
	flag.StringVar(&key, "ssh-key", "", "The location of private key")

	flag.Parse()
	os.Exit(m.Run())
}
