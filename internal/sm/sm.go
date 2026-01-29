package sm

import "fmt"

// Provider constants
const (
	ProviderVault = "vault"
	// ProviderAWS = "aws" // Future implementation
)

// KeyValue represents secret key-value data.
type KeyValue struct {
	PrivateKey        []byte
	PublicKey         []byte
	RequirePassphrase bool
	Comment           string
}

// Get retrieves key-value data from a provider at the given path
func Get(provider, path string) (*KeyValue, error) {
	switch provider {
	case ProviderVault:
		client, err := NewVaultClient()
		if err != nil {
			return nil, err
		}
		return client.GetKV(path)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// Store stores key-value data to a provider at the given path
func Store(provider, path string, kv *KeyValue) error {
	switch provider {
	case ProviderVault:
		client, err := NewVaultClient()
		if err != nil {
			return err
		}
		return client.StoreKV(path, kv)
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}
