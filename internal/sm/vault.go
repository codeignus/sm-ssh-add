package sm

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/vault/api/auth/approle"

	vaultapi "github.com/hashicorp/vault/api"
)

// VaultApproleConfig is the interface for Vault Approle authentication configuration.
// Using an interface avoids circular imports with the config package.
type VaultApproleConfig interface {
	GetVaultApproleRoleID() string
}

// VaultClient implements SecretManager interface for HashiCorp Vault KV v2
type VaultClient struct {
	client *vaultapi.Client
}

// getVaultAddress returns the Vault/OpenBao address from environment variables.
func getVaultAddress() string {
	addr := os.Getenv("BAO_ADDR")
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	return addr
}

// getVaultToken returns the Vault/OpenBao token from environment variables.
func getVaultToken() string {
	token := os.Getenv("BAO_TOKEN")
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	return token
}

// promptForSecretID prompts the user to enter their AppRole Secret ID.
func promptForSecretID() (string, error) {
	fmt.Fprintln(os.Stderr, "Generate a single-use Secret ID using command similar to below:")
	fmt.Fprintln(os.Stderr, "vault write -f auth/approle/role/sm-ssh-add/secret-id")
	fmt.Fprint(os.Stderr, "Enter Vault/OpenBao AppRole Secret ID: ")

	var secretID string
	_, err := fmt.Scanln(&secretID)
	if err != nil {
		return "", wrapError(err, "failed to read secret ID")
	}
	if secretID == "" {
		return "", fmt.Errorf("secret ID cannot be empty")
	}
	return secretID, nil
}

// appRoleLogin performs AppRole authentication and sets the token on the client.
func appRoleLogin(client *vaultapi.Client, roleID string) error {
	secretID := os.Getenv("VAULT_APPROLE_SECRET_ID")
	if secretID == "" {
		var err error
		secretID, err = promptForSecretID()
		if err != nil {
			return err
		}
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		roleID,
		&approle.SecretID{FromString: secretID},
	)
	if err != nil {
		return wrapError(err, "failed to initialize AppRole auth")
	}

	_, err = client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return wrapError(err, "failed to login with AppRole")
	}

	return nil
}

// NewVaultClient creates a new Vault client using environment variables.
// If cfg is provided and contains VaultApproleRoleID, performs Approle login instead of using token.
func NewVaultClient(cfg VaultApproleConfig) (*VaultClient, error) {
	addr := getVaultAddress()
	if addr == "" {
		return nil, fmt.Errorf("vault address required: set BAO_ADDR or VAULT_ADDR")
	}

	config := vaultapi.DefaultConfig()
	// Explicitly set address from env to preserve scheme (http:// vs https://)
	config.Address = addr

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, wrapError(err, "failed to create vault client")
	}

	// Check if config has VaultApproleRoleID field set
	var roleID string
	if cfg != nil {
		roleID = cfg.GetVaultApproleRoleID()
	}

	// Authenticate: AppRole if configured, otherwise token
	if roleID != "" {
		if err := appRoleLogin(client, roleID); err != nil {
			return nil, err
		}
	} else {
		token := getVaultToken()
		if token == "" {
			return nil, fmt.Errorf("vault token required: set BAO_TOKEN or VAULT_TOKEN")
		}
		client.SetToken(token)
	}

	// Verify connection
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		return nil, wrapError(err, ErrVaultConnection.Error())
	}

	return &VaultClient{client: client}, nil
}

// Get retrieves key-value data from Vault KV v2 at the given path
func (v *VaultClient) Get(path string) (*KeyValue, error) {
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
	requirePassphraseStr, ok := data["require_passphrase"].(string)
	if ok {
		var err error
		requirePassphrase, err = strconv.ParseBool(requirePassphraseStr)
		if err != nil {
			return nil, wrapError(err, "failed to parse require_passphrase")
		}
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

// Store stores key-value data in Vault KV v2 at the given path
func (v *VaultClient) Store(path string, kv *KeyValue) error {
	secretData := map[string]interface{}{
		"private_key":        string(kv.PrivateKey),
		"public_key":         string(kv.PublicKey),
		"require_passphrase": fmt.Sprintf("%v", kv.RequirePassphrase),
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

// CheckExists checks if a key already exists at the given path
func (v *VaultClient) CheckExists(path string) (bool, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return false, wrapError(err, "failed to check path existence")
	}

	// KV v2: If secret.Data is nil or empty, path doesn't exist
	if secret == nil || secret.Data == nil {
		return false, nil
	}

	data, ok := secret.Data["data"]
	if !ok || data == nil {
		return false, nil
	}

	// If data map exists and is not empty, path exists
	if dataMap, ok := data.(map[string]interface{}); ok {
		return len(dataMap) > 0, nil
	}

	return false, nil
}
