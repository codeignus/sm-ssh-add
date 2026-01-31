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

// Get retrieves key-value data from a provider at the given path
func Get(cfg *config.Config, path string) (*KeyValue, error) {
	provider := cfg.DefaultProvider
	switch provider {
	case config.ProviderVault:
		client, err := NewVaultClient(cfg)
		if err != nil {
			return nil, err
		}
		return client.GetKV(path)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// Store stores key-value data to a provider at the given path
func Store(cfg *config.Config, path string, kv *KeyValue) error {
	provider := cfg.DefaultProvider
	switch provider {
	case config.ProviderVault:
		client, err := NewVaultClient(cfg)
		if err != nil {
			return err
		}
		return client.StoreKV(path, kv)
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}
