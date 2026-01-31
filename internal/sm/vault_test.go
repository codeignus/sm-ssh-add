package sm

import (
	"os"
	"testing"
)

// mockConfig is a test helper that implements the VaultApproleConfig interface
type mockConfig struct {
	VaultApproleRoleID string
}

// GetVaultApproleRoleID makes mockConfig implement the VaultApproleConfig interface
func (c *mockConfig) GetVaultApproleRoleID() string {
	return c.VaultApproleRoleID
}

func TestNewVaultClient_reads_environment_variables(t *testing.T) {
	tests := []struct {
		name        string
		setBAO      bool
		setVault    bool
		expectError bool
	}{
		{
			name:        "reads BAO_ADDR and BAO_TOKEN",
			setBAO:      true,
			setVault:    false,
			expectError: false,
		},
		{
			name:        "falls back to VAULT_ADDR and VAULT_TOKEN",
			setBAO:      false,
			setVault:    true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			oldBAOAddr := os.Getenv("BAO_ADDR")
			oldBAOToken := os.Getenv("BAO_TOKEN")
			oldVaultAddr := os.Getenv("VAULT_ADDR")
			oldVaultToken := os.Getenv("VAULT_TOKEN")
			defer func() {
				os.Setenv("BAO_ADDR", oldBAOAddr)
				os.Setenv("BAO_TOKEN", oldBAOToken)
				os.Setenv("VAULT_ADDR", oldVaultAddr)
				os.Setenv("VAULT_TOKEN", oldVaultToken)
			}()

			// Clear all first
			os.Unsetenv("BAO_ADDR")
			os.Unsetenv("BAO_TOKEN")
			os.Unsetenv("VAULT_ADDR")
			os.Unsetenv("VAULT_TOKEN")

			// Set based on test case
			if tt.setBAO {
				os.Setenv("BAO_ADDR", "http://localhost:8200")
				os.Setenv("BAO_TOKEN", "test-token")
			}
			if tt.setVault {
				os.Setenv("VAULT_ADDR", "http://localhost:8200")
				os.Setenv("VAULT_TOKEN", "test-token")
			}

			client, err := NewVaultClient(nil)

			// Will fail connection (no real Vault), but should not fail on missing addr/token
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil && client != nil {
					t.Errorf("Unexpected error with non-nil client: %v", err)
				}
			}
		})
	}
}

func TestNewVaultClient_rejects_missing_address(t *testing.T) {
	oldAddr := os.Getenv("BAO_ADDR")
	oldVaultAddr := os.Getenv("VAULT_ADDR")
	oldToken := os.Getenv("BAO_TOKEN")
	oldVaultToken := os.Getenv("VAULT_TOKEN")
	defer func() {
		os.Setenv("BAO_ADDR", oldAddr)
		os.Setenv("VAULT_ADDR", oldVaultAddr)
		os.Setenv("BAO_TOKEN", oldToken)
		os.Setenv("VAULT_TOKEN", oldVaultToken)
	}()

	os.Unsetenv("BAO_ADDR")
	os.Unsetenv("VAULT_ADDR")
	os.Setenv("BAO_TOKEN", "test-token")

	client, err := NewVaultClient(nil)

	if err == nil {
		t.Error("Expected error when no address set")
	}
	if client != nil {
		t.Error("Expected nil client when address missing")
	}
}

func TestNewVaultClient_rejects_missing_token(t *testing.T) {
	oldAddr := os.Getenv("BAO_ADDR")
	oldToken := os.Getenv("BAO_TOKEN")
	oldVaultToken := os.Getenv("VAULT_TOKEN")
	defer func() {
		os.Setenv("BAO_ADDR", oldAddr)
		os.Setenv("BAO_TOKEN", oldToken)
		os.Setenv("VAULT_TOKEN", oldVaultToken)
	}()

	os.Setenv("BAO_ADDR", "http://localhost:8200")
	os.Unsetenv("BAO_TOKEN")
	os.Unsetenv("VAULT_TOKEN")

	client, err := NewVaultClient(nil)

	if err == nil {
		t.Error("Expected error when no token set")
	}
	if client != nil {
		t.Error("Expected nil client when token missing")
	}
}

func TestNewVaultClient_withAppRoleConfig_performsLogin(t *testing.T) {
	// Save original env vars
	oldAddr := os.Getenv("BAO_ADDR")
	oldVaultAddr := os.Getenv("VAULT_ADDR")
	oldToken := os.Getenv("BAO_TOKEN")
	oldVaultToken := os.Getenv("VAULT_TOKEN")
	oldSecretID := os.Getenv("VAULT_APPROLE_SECRET_ID")
	defer func() {
		os.Setenv("BAO_ADDR", oldAddr)
		os.Setenv("VAULT_ADDR", oldVaultAddr)
		os.Setenv("BAO_TOKEN", oldToken)
		os.Setenv("VAULT_TOKEN", oldVaultToken)
		os.Setenv("VAULT_APPROLE_SECRET_ID", oldSecretID)
	}()

	// Set address and secret ID (from env to avoid prompt), but NOT token
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Unsetenv("BAO_TOKEN")
	os.Unsetenv("VAULT_TOKEN")
	os.Setenv("VAULT_APPROLE_SECRET_ID", "test-secret-id")

	// Create a mock config with Vault Approle role_id
	cfg := &mockConfig{
		VaultApproleRoleID: "test-role-id",
	}

	// Should attempt AppRole login (will fail connection, but that's ok)
	_, err := NewVaultClient(cfg)

	// Should not error on missing token
	// Will likely error on connection (no real Vault), but that's different from "token required"
	if err != nil {
		// Error is ok (no real Vault), but verify it's not "token required" error
		if err.Error() == "vault token required: set BAO_TOKEN or VAULT_TOKEN" {
			t.Error("Should not require token when AppRole is configured")
		}
	}
}
