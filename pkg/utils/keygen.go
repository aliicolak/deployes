package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// GenerateSSHKeyPair generates a 4096-bit RSA key pair.
// Returns privateKeyPEM and publicKey (authorized_keys format).
func GenerateSSHKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	// Encode Private Key to PEM
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privatePEM := string(pem.EncodeToMemory(&privBlock))

	// Encode Public Key to Authorized Key format
	publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	publicBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	publicStr := string(publicBytes)

	return privatePEM, publicStr, nil
}
