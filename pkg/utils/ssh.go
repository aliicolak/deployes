package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHConnectionMode defines the host key verification mode
type SSHConnectionMode int

const (
	// SSHModeStrict requires the host key to already be in known_hosts.
	// An unknown host is rejected.
	SSHModeStrict SSHConnectionMode = iota
	// SSHModeTrustOnFirstUse accepts an unknown host key, persists it to
	// known_hosts, and verifies against it on every later connection.
	SSHModeTrustOnFirstUse
)

// knownHostsMu serializes read-modify-write cycles on the known_hosts file.
// The deploy worker, the terminal handler and the server service can all open
// connections concurrently.
var knownHostsMu sync.Mutex

// CreateSSHClient creates SSH connection with private key or password.
// Uses trust-on-first-use host key verification: the first connection to a host
// records its key, and any later key change is rejected.
func CreateSSHClient(host string, port int, username string, secret string) (*ssh.Client, error) {
	return CreateSSHClientWithMode(host, port, username, secret, SSHModeTrustOnFirstUse)
}

// CreateSSHClientWithMode creates SSH connection with specified host key verification mode
func CreateSSHClientWithMode(host string, port int, username string, secret string, mode SSHConnectionMode) (*ssh.Client, error) {

	var authMethods []ssh.AuthMethod

	// Try to parse as private key
	signer, err := ssh.ParsePrivateKey([]byte(secret))
	if err == nil {
		// It's a valid private key
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		// Not a private key, treat as password
		authMethods = append(authMethods, ssh.Password(secret))
	}

	// Get host key callback based on mode
	hostKeyCallback, err := getHostKeyCallback(mode)
	if err != nil {
		return nil, fmt.Errorf("failed to create host key callback: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			return nil, fmt.Errorf("failed to connect to server: %w", err)
		}
		return nil, err
	}

	return conn, nil
}

// getHostKeyCallback returns a callback that verifies the host key against
// known_hosts. In trust-on-first-use mode an unknown host is learned and
// persisted; in either mode a key that contradicts a stored one is rejected.
func getHostKeyCallback(mode SSHConnectionMode) (ssh.HostKeyCallback, error) {
	path := knownHostsPath()

	if mode != SSHModeStrict && mode != SSHModeTrustOnFirstUse {
		return nil, fmt.Errorf("unknown SSH connection mode")
	}

	if mode == SSHModeStrict {
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("known_hosts file not found at %s - connect once in trust-on-first-use mode or provision it manually", path)
		}
	} else if err := ensureKnownHostsFile(path); err != nil {
		return nil, err
	}

	verify, err := knownhosts.New(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := verify(hostname, remote, key)
		if err == nil {
			return nil
		}

		var keyErr *knownhosts.KeyError
		if !errors.As(err, &keyErr) {
			return err
		}

		// Want is non-empty when known_hosts holds a *different* key for this
		// host. That is the man-in-the-middle case and is never acceptable.
		if len(keyErr.Want) > 0 {
			return fmt.Errorf(
				"host key mismatch for %s: the server presented %s but %s records a different key - "+
					"this may be a man-in-the-middle attack; if the host was legitimately rebuilt, remove its entry from %s",
				hostname, ssh.FingerprintSHA256(key), path, path,
			)
		}

		// Host is not in known_hosts at all.
		if mode == SSHModeStrict {
			return fmt.Errorf("unknown host %s (fingerprint %s) and strict mode is enabled", hostname, ssh.FingerprintSHA256(key))
		}

		if err := appendKnownHost(path, hostname, key); err != nil {
			return fmt.Errorf("failed to record host key for %s: %w", hostname, err)
		}
		log.Printf("🔑 SSH: learned host key for %s (fingerprint: %s)", hostname, ssh.FingerprintSHA256(key))
		return nil
	}, nil
}

// ensureKnownHostsFile creates the known_hosts file and its directory when
// missing, so that the first ever connection has something to read.
func ensureKnownHostsFile(path string) error {
	knownHostsMu.Lock()
	defer knownHostsMu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create %s: %w", filepath.Dir(path), err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	return f.Close()
}

// appendKnownHost records a newly seen host key.
func appendKnownHost(path, hostname string, key ssh.PublicKey) error {
	knownHostsMu.Lock()
	defer knownHostsMu.Unlock()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	line := knownhosts.Line([]string{knownhosts.Normalize(hostname)}, key)
	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil
}

// knownHostsPath returns the path to the known_hosts file
func knownHostsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".ssh/known_hosts"
	}
	return filepath.Join(homeDir, ".ssh", "known_hosts")
}
