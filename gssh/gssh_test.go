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
	t.Run("TestPassAuth", secure(t, authTest))
	t.Run("Ping", func(t *testing.T) {
		err := Ping(addr, user, passwd, key)
		assert.NoError(t, err)
	})
	t.Run("Ping error", func(t *testing.T) {
		err := Ping(addr, user, "error", key)
		if err != nil {
			assert.Contains(t, err.Error(), "ssh: unable to authenticate")
		}
	})
}

func TestGSSHInsecure(t *testing.T) {
	t.Run("TestPassAuth", insecure(t, authTest))
	t.Run("TestCmdOutPipe", insecure(t, outPipeTest))
	t.Run("TestSetEnv", insecure(t, envTest))
	t.Run("TestClientCmd", insecure(t, cliCmdTest))
}

func insecure(t *testing.T, callback func(t *testing.T, cli *Client)) func(t *testing.T) {
	opts := NewOptions()
	opts.Username = user
	opts.Password = passwd
	opts.Addr = addr
	opts.Key = key

	cli, err := NewInsecure(opts)
	assert.NoError(t, err)
	err = cli.Dial()
	assert.NoError(t, err)
	return func(t *testing.T) {
		callback(t, cli)
	}
}

func secure(t *testing.T, callback func(t *testing.T, cli *Client)) func(t *testing.T) {
	opts := NewOptions()
	opts.Username = user
	opts.Password = passwd
	opts.Addr = addr
	opts.Key = key

	cli, err := New(opts)
	assert.NoError(t, err)
	cli.WithHostKeyCallback(AutoFixedHostKeyCallback)
	err = cli.Dial()
	assert.NoError(t, err)
	return func(t *testing.T) {
		defer func() { _ = cli.Close() }()
		callback(t, cli)
	}
}

func cliCmdTest(t *testing.T, cli *Client) {
	output, err := cli.CombinedOutput("echo \"Hello, world!\"")
	assert.NoError(t, err)

	assert.Equal(t, "Hello, world!\n", string(output))
}

func authTest(t *testing.T, cli *Client) {
	cmd, err := cli.Command("echo", "Hello, world!")
	assert.NoError(t, err)

	output, err := cmd.Output()
	assert.NoError(t, err)

	assert.Equal(t, "Hello, world!\n", string(output))
}

func outPipeTest(t *testing.T, cli *Client) {
	cmd, err := cli.Command("n=1;while [ $n -le 4 ];do echo $n;((n++));done")
	assert.NoError(t, err)
	var lines []string
	err = cmd.OutputPipe(func(line []byte) error {
		// t.Log(line)
		lines = append(lines, string(line))
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3", "4"}, lines)
}

func envTest(t *testing.T, cli *Client) {
	cmd, err := cli.Command("echo", "Hello, $TEST_ENV_NAME!")
	assert.NoError(t, err)
	err = cmd.Setenv([]string{"TEST_ENV_NAME=GSSH"})
	assert.NoError(t, err)
	output, err := cmd.Output()
	assert.NoError(t, err)

	assert.Equal(t, "Hello, GSSH!\n", string(output))
}

func TestMain(m *testing.M) {
	flag.StringVar(&addr, "addr", "", "The host of ssh")
	flag.StringVar(&user, "user", "root", "The username of client")
	flag.StringVar(&passwd, "pass", "", "The password of user")
	flag.StringVar(&key, "ssh-key", "", "The location of private key")

	flag.Parse()
	if addr != "" {
		os.Exit(m.Run())
	}
}
