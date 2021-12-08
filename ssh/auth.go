package ssh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)


// Auth returns a single auth method.
func Auth(opts *Options) ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod
	if opts.Password != "" {
		methods = append(methods, ssh.Password(opts.Password))
	}
	if opts.Key != "" {
		signer, err := GetSigner(opts.Key, opts.Passphrase)
		if err != nil {
			return nil, err
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}
	return methods, errors.New("no auth method")
}

// HasAgent checks if ssh agent exists.
func HasAgent() bool {
	return os.Getenv("SSH_AUTH_SOCK") != ""
}

// AgentAuth auth via ssh agent, (Unix systems only)
func AgentAuth() (ssh.AuthMethod, error) {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, fmt.Errorf("could not find ssh agent: %w", err)
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
}

// GetSigner returns ssh.Signer from private key file.
func GetSigner(sshkey string, passphrase string) (signer ssh.Signer, err error) {
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
