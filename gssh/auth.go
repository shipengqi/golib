package gssh

import (
	"errors"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Auth returns a single ssh.AuthMethod.
func Auth(opts *Options) (ssh.AuthMethod, error) {
	var (
		auth ssh.AuthMethod
		err  error
	)
	if opts.UseAgent && HasAgent() {
		if auth, err = Agent(); err == nil {
			return auth, nil
		}
	}
	if opts.Key != "" {
		if auth, err = Key(opts.Key, opts.Passphrase); err == nil {
			return auth, nil
		}
	}
	if opts.Password != "" {
		auth = Password(opts.Password)
		return auth, nil
	}
	return nil, errors.New("no auth method")
}

// HasAgent checks if ssh agent exists.
func HasAgent() bool {
	return os.Getenv("SSH_AUTH_SOCK") != ""
}

// Agent returns ssh.AuthMethod of ssh agent, (Unix systems only).
func Agent() (ssh.AuthMethod, error) {
	if !HasAgent() {
		return nil, errors.New("no agent")
	}
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
}

// Password returns ssh.AuthMethod of password.
func Password(pass string) ssh.AuthMethod {
	return ssh.Password(pass)
}

// Key returns ssh.AuthMethod from private key file.
func Key(sshkey string, passphrase string) (ssh.AuthMethod, error) {
	signer, err := GetSigner(sshkey, passphrase)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

// GetSigner returns ssh.Signer from private key file.
func GetSigner(sshkey, passphrase string) (signer ssh.Signer, err error) {
	data, err := ioutil.ReadFile(sshkey)
	if err != nil {
		return
	}
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(data, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(data)
	}
	return
}
