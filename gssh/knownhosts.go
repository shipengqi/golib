package gssh

import (
	"errors"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func AutoFixedHostKeyCallback(host string, remote net.Addr, key ssh.PublicKey) error {
	found, err := VerifyKnownHost("", host, remote, key)
	if found && err != nil {
		return err
	}
	if found && err == nil {
		return nil
	}
	// Add the new host to known hosts file.
	return AppendKnownHost("", host, remote, key)
}

// DefaultHostKeyCallback returns host key callback from default known_hosts file.
func DefaultHostKeyCallback() (ssh.HostKeyCallback, error) {
	fpath, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, err
	}
	return knownhosts.New(fpath)
}

// DefaultKnownHostsPath returns the path of default knows_hosts file.
func DefaultKnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ssh/known_hosts"), err
}

// VerifyKnownHost reports whether the given host in known hosts file and valid.
func VerifyKnownHost(fpath, host string, remote net.Addr, key ssh.PublicKey) (bool, error) {
	var (
		keyErr   *knownhosts.KeyError
		callback ssh.HostKeyCallback
		err      error
	)
	if fpath != "" {
		callback, err = knownhosts.New(fpath)
	} else {
		callback, err = DefaultHostKeyCallback()
	}
	if err != nil {
		return false, err
	}

	// check if host already exists
	if err = callback(host, remote, key); err == nil {
		return true, nil
	} else if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		// The known_hosts file contains the given host, but with different key.
		return true, keyErr
	}
	return false, nil
}

// AppendKnownHost appends a host to known hosts file.
func AppendKnownHost(fpath, host string, remote net.Addr, key ssh.PublicKey) error {
	if fpath == "" {
		f, err := DefaultKnownHostsPath()
		if err != nil {
			return err
		}
		fpath = f
	}

	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	remoteNormalized := knownhosts.Normalize(remote.String())
	hostNormalized := knownhosts.Normalize(host)
	addresses := []string{remoteNormalized}
	if hostNormalized != remoteNormalized {
		addresses = append(addresses, hostNormalized)
	}

	_, err = f.WriteString(knownhosts.Line(addresses, key) + "\n")

	return err
}
