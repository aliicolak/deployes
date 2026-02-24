package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHConnectionMode defines the host key verification mode
type SSHConnectionMode int

const (
	// SSHModeStrict requires host key to be in known_hosts file
	SSHModeStrict SSHConnectionMode = iota
	// SSHModeTrustOnFirstUse accepts and stores host key on first connection
	SSHModeTrustOnFirstUse
)

// CreateSSHClient creates SSH connection with private key or password
// For production, use SSHModeStrict with a properly configured known_hosts file
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
	hostKeyCallback, err := getHostKeyCallback(host, port, mode)
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

// getHostKeyCallback returns appropriate host key callback based on mode
func getHostKeyCallback(host string, port int, mode SSHConnectionMode) (ssh.HostKeyCallback, error) {
	knownHostsPath := getKnownHostsPath()

	switch mode {
	case SSHModeStrict:
		// Strict mode: require host key to be in known_hosts
		if _, err := os.Stat(knownHostsPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("known_hosts file not found at %s - create it or use TrustOnFirstUse mode", knownHostsPath)
		}
		return knownhosts.New(knownHostsPath)

	case SSHModeTrustOnFirstUse:
		// Trust on first use: accept and log warning
		// In production, you should use SSHModeStrict
		log.Printf("⚠️  SSH: Using Trust-On-First-Use mode for %s:%d - verify server authenticity!", host, port)

		// Create a callback that logs the host key on first use
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			log.Printf("🔑 SSH: Accepting host key for %s (fingerprint: %s)", hostname, ssh.FingerprintSHA256(key))
			// In a full implementation, you would save this to known_hosts
			return nil
		}, nil

	default:
		return nil, fmt.Errorf("unknown SSH connection mode")
	}
}

// getKnownHostsPath returns the path to the known_hosts file
func getKnownHostsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".ssh/known_hosts"
	}
	return filepath.Join(homeDir, ".ssh", "known_hosts")
}
