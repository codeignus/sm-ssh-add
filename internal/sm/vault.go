package sm

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// VaultClient implements SecretManager interface for HashiCorp Vault KV v2
type VaultClient struct {
	client *vaultapi.Client
}

// NewVaultClient creates a new Vault client using environment variables.
// Reads BAO_ADDR/BAO_TOKEN first, falls back to VAULT_ADDR/VAULT_TOKEN for compatibility.
func NewVaultClient() (*VaultClient, error) {
	// Check environment variables (OpenBao first, then Vault for compatibility)
	addr := os.Getenv("BAO_ADDR")
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}

	token := os.Getenv("BAO_TOKEN")
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	if addr == "" {
		return nil, fmt.Errorf("vault address required: set BAO_ADDR or VAULT_ADDR")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token required: set BAO_TOKEN or VAULT_TOKEN")
	}

	config := vaultapi.DefaultConfig()
	// Explicitly set address from env to preserve scheme (http:// vs https://)
	config.Address = addr

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, wrapError(err, "failed to create vault client")
	}

	client.SetToken(token)

	// Verify connection by making a simple request
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		return nil, wrapError(err, ErrVaultConnection.Error())
	}

	return &VaultClient{
		client: client,
	}, nil
}

// GetKV retrieves key-value data from Vault KV v2 at the given path
func (v *VaultClient) GetKV(path string) (*KeyValue, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, wrapError(err, "failed to read from vault")
	}

	if secret == nil {
		return nil, ErrPathNotFound
	}

	// Handle KV v2 data wrapper
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidKeyFormat
	}

	privateKey, ok := data["private_key"].(string)
	if !ok || privateKey == "" {
		return nil, ErrInvalidKeyFormat
	}

	publicKey, ok := data["public_key"].(string)
	if !ok || publicKey == "" {
		return nil, ErrInvalidKeyFormat
	}

	requirePassphrase := false
	if rp, ok := data["require_passphrase"].(bool); ok {
		requirePassphrase = rp
	}

	comment := ""
	if c, ok := data["comment"].(string); ok {
		comment = c
	}

	return &KeyValue{
		PrivateKey:        []byte(privateKey),
		PublicKey:         []byte(publicKey),
		RequirePassphrase: requirePassphrase,
		Comment:           comment,
	}, nil
}

// StoreKV stores key-value data in Vault KV v2 at the given path
func (v *VaultClient) StoreKV(path string, kv *KeyValue) error {
	secretData := map[string]interface{}{
		"private_key":        string(kv.PrivateKey),
		"public_key":         string(kv.PublicKey),
		"require_passphrase": kv.RequirePassphrase,
		"comment":            kv.Comment,
	}

	data := map[string]interface{}{
		"data": secretData,
	}

	_, err := v.client.Logical().Write(path, data)
	if err != nil {
		return wrapError(err, "failed to write to vault")
	}

	return nil
}
