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

// InitProvider creates and initializes a Provider based on the config
// The provider client is created once here and reused for all operations
func InitProvider(cfg *config.Config) (Provider, error) {
	switch cfg.DefaultProvider {
	case config.ProviderVault:
		return NewVaultClient(cfg)
	// Future: add case for AWS, Azure, etc.
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
	}
}
