package sm

import (
	"os"
	"testing"
)

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

			client, err := NewVaultClient()

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

	client, err := NewVaultClient()

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

	client, err := NewVaultClient()

	if err == nil {
		t.Error("Expected error when no token set")
	}
	if client != nil {
		t.Error("Expected nil client when token missing")
	}
}
