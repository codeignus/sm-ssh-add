package ssh

import (
	"net"
	"os"

	"github.com/codeignus/sm-ssh-add/internal/config"
	"github.com/codeignus/sm-ssh-add/internal/sm"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Agent handles SSH agent operations and provides methods to add keys and list existing keys
type Agent struct {
	client agent.ExtendedAgent
	conn   net.Conn
}

// NewAgent connects to the SSH agent using the socket path from environment
func NewAgent(cfg *config.Config) (*Agent, error) {
	socketPath := os.Getenv("SSH_AUTH_SOCK")
	if socketPath == "" {
		return nil, sm.ErrSSHAgentNotFound
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, wrapError(err, "failed to connect to ssh-agent")
	}

	agentClient := agent.NewClient(conn)

	return &Agent{
		client: agentClient,
		conn:   conn,
	}, nil
}

// AddKey adds a key pair to the SSH agent after checking if it already exists
func (a *Agent) AddKey(keyPair *KeyPair) error {
	var privateKey interface{}
	var signer ssh.Signer
	var err error

	// Parse private key with or without passphrase
	if keyPair.Passphrase != nil {
		privateKey, err = ssh.ParseRawPrivateKeyWithPassphrase(keyPair.PrivateKey, []byte(*keyPair.Passphrase))
		if err != nil {
			return wrapError(err, "failed to parse private key with passphrase")
		}
		// Also create a signer for fingerprint checking
		signer, err = ssh.NewSignerFromKey(privateKey)
		if err != nil {
			return wrapError(err, "failed to create signer from private key")
		}
	} else {
		privateKey, err = ssh.ParseRawPrivateKey(keyPair.PrivateKey)
		if err != nil {
			return wrapError(err, "failed to parse private key")
		}
		// Also create a signer for fingerprint checking
		signer, err = ssh.NewSignerFromKey(privateKey)
		if err != nil {
			return wrapError(err, "failed to create signer from private key")
		}
	}

	keys, err := a.client.List()
	if err != nil {
		return wrapError(err, "failed to list keys in agent")
	}

	existingFingerprint := ssh.FingerprintSHA256(signer.PublicKey())
	for _, key := range keys {
		if ssh.FingerprintSHA256(key) == existingFingerprint {
			return sm.ErrKeyExistsInAgent
		}
	}

	addedKey := agent.AddedKey{
		PrivateKey:       privateKey,
		Comment:          keyPair.Comment,
		LifetimeSecs:     0,
		ConfirmBeforeUse: false,
	}

	err = a.client.Add(addedKey)
	if err != nil {
		return wrapError(err, "failed to add key to agent")
	}

	return nil
}

// List returns all keys currently loaded in the SSH agent
func (a *Agent) List() ([]*agent.Key, error) {
	keys, err := a.client.List()
	if err != nil {
		return nil, wrapError(err, "failed to list keys")
	}
	return keys, nil
}

// Close closes the connection to the SSH agent
func (a *Agent) Close() error {
	if a.conn != nil {
		err := a.conn.Close()
		a.client = nil
		a.conn = nil
		return err
	}
	return nil
}

// KeyExists checks if a key with the given fingerprint is already loaded in the agent
func (a *Agent) KeyExists(fingerprint string) (bool, error) {
	keys, err := a.client.List()
	if err != nil {
		return false, wrapError(err, "failed to list keys")
	}

	for _, key := range keys {
		if ssh.FingerprintSHA256(key) == fingerprint {
			return true, nil
		}
	}

	return false, nil
}
