package utils

import (
	"bytes"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

func RunSSHCommand(client *ssh.Client, command string) (string, error) {

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)

	output := stdout.String() + stderr.String()
	return output, err
}

// RunLocalCommand executes a command on the local machine
func RunLocalCommand(command string) (string, error) {
	// Split command into parts for exec
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()
	return output, err
}
