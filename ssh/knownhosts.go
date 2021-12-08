package ssh

import (
	"errors"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// DefaultHostKeyCallback returns host key callback from default known_hosts file.
func DefaultHostKeyCallback() (ssh.HostKeyCallback, error) {
	fpath, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, err
	}
	return knownhosts.New(fpath)
}

// CheckKnownHost reports whether the given host in known_hosts file.
func CheckKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownHosts string) (found bool, err error) {
	var keyErr *knownhosts.KeyError

	// Fallback to default known_hosts file
	if knownHosts == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return false, err
		}

		knownHosts = path
	}

	// Get host key callback
	callback, err := knownhosts.New(knownHosts)

	if err != nil {
		return false, err
	}

	// check if host already exists.
	err = callback(host, remote, key)

	// Known host already exists.
	if err == nil {
		return true, nil
	}

	// Make sure that the error returned from the callback is host not in file error.
	// If keyErr.Want is greater than 0 length, that means host is in file with different key.
	if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		return true, keyErr
	}

	// Some other error occurred and safest way to handle is to pass it back to user.
	if err != nil {
		return false, err
	}

	// Key is not trusted because it is not in the file.
	return false, nil
}

// AppendKnownHost appends a host to known hosts file.
func AppendKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownHosts string) (err error) {
	// Fallback to default known_hosts file
	if knownHosts == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return err
		}
		knownHosts = path
	}

	f, err := os.OpenFile(knownHosts, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
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

// DefaultKnownHostsPath returns the path of default knows_hosts file.
func DefaultKnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ssh/known_hosts"), err
}
