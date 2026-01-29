package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// KeyPair represents an SSH key pair with private and public components
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
	Comment    string
}

// GenerateKeyPair generates a new ed25519 SSH key pair and marshals it to OpenSSH format
// If passphrase is non-empty, the private key will be encrypted using the provided passphrase
func GenerateKeyPair(comment string, passphrase []byte) (*KeyPair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, wrapError(err, "failed to generate ed25519 key")
	}

	var privateKeyPEM []byte
	if len(passphrase) > 0 {
		privateKeyBlock, err := ssh.MarshalPrivateKeyWithPassphrase(privKey, comment, passphrase)
		if err != nil {
			return nil, wrapError(err, "failed to marshal private key with passphrase")
		}
		privateKeyPEM = pem.EncodeToMemory(privateKeyBlock)
	} else {
		privateKeyBlock, err := ssh.MarshalPrivateKey(privKey, comment)
		if err != nil {
			return nil, wrapError(err, "failed to marshal private key")
		}
		privateKeyPEM = pem.EncodeToMemory(privateKeyBlock)
	}

	pubKeySSH, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, wrapError(err, "failed to convert public key")
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(pubKeySSH)

	return &KeyPair{
		PrivateKey: privateKeyPEM,
		PublicKey:  publicKeyBytes,
		Comment:    comment,
	}, nil
}
