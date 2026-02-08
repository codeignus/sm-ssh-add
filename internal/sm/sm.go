package sm

import (
	"fmt"

	"github.com/codeignus/sm-ssh-add/internal/config"
)

// KeyValue represents secret key-value data.
type KeyValue struct {
	PrivateKey        []byte
	PublicKey         []byte
	RequirePassphrase bool
	Comment           string
}

// Provider defines the interface for secret manager providers
type Provider interface {
	Get(path string) (*KeyValue, error)
	Store(path string, kv *KeyValue) error
	CheckExists(path string) (bool, error)
}

// providerImpl implements the Provider interface
// It holds provider clients
type providerImpl struct {
	vaultClient *VaultClient
	// Future: add awsClient, azureClient, etc.
}

// Get retrieves key-value data from a provider at the given path
// This method checks which provider is initialized and delegates accordingly
func (p *providerImpl) Get(path string) (*KeyValue, error) {
	switch {
	case p.vaultClient != nil:
		// Vault-specific: GetKV returns KeyValue directly
		return p.vaultClient.GetKV(path)
	// Future: add case for AWS, Azure, etc.
	default:
		return nil, fmt.Errorf("no provider initialized")
	}
}

// Store stores key-value data to a provider at the given path
func (p *providerImpl) Store(path string, kv *KeyValue) error {
	switch {
	case p.vaultClient != nil:
		return p.vaultClient.StoreKV(path, kv)
	default:
		return fmt.Errorf("no provider initialized")
	}
}

// CheckExists checks if a key already exists at the given path
func (p *providerImpl) CheckExists(path string) (bool, error) {
	switch {
	case p.vaultClient != nil:
		return p.vaultClient.CheckExists(path)
	default:
		return false, fmt.Errorf("no provider initialized")
	}
}

// InitProvider creates and initializes a Provider based on the config
// The provider client is created once here and reused for all operations
func InitProvider(cfg *config.Config) (Provider, error) {
	impl := &providerImpl{}

	switch cfg.DefaultProvider {
	case config.ProviderVault:
		client, err := NewVaultClient(cfg)
		if err != nil {
			return nil, err
		}
		impl.vaultClient = client
	// Future: add case for AWS, Azure, etc.
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
	}

	return impl, nil
}
