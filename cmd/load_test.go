package cmd

import (
	"testing"

	"github.com/codeignus/sm-ssh-add/internal/config"
	"github.com/codeignus/sm-ssh-add/internal/sm"
)

// mockProviderForLoad is a mock implementation of sm.Provider for testing
type mockProviderForLoad struct{}

func (m *mockProviderForLoad) Get(path string) (*sm.KeyValue, error) {
	// Return mock key data for testing
	return &sm.KeyValue{
		PrivateKey:        []byte("mock-private-key"),
		PublicKey:         []byte("mock-public-key"),
		RequirePassphrase: false,
		Comment:           "test",
	}, nil
}

func (m *mockProviderForLoad) Store(path string, kv *sm.KeyValue) error {
	return nil
}

func (m *mockProviderForLoad) CheckExists(path string) (bool, error) {
	return false, nil
}

// TestLoadFromConfig_EmptyPaths tests --from-config with empty paths
func TestLoadFromConfig_EmptyPaths(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{},
	}
	provider := &mockProviderForLoad{}

	err := Load(provider, cfg, []string{"--from-config"})

	if err == nil {
		t.Error("expected error for empty paths, got nil")
	}
	if err != nil && err.Error() != "no paths configured" {
		t.Errorf("expected 'no paths configured' error, got: %v", err)
	}
}

// TestLoadNoArguments tests error when no arguments provided
func TestLoadNoArguments(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{"secret/ssh/test"},
	}
	provider := &mockProviderForLoad{}

	err := Load(provider, cfg, []string{})

	if err == nil {
		t.Error("expected error for no arguments, got nil")
	}
}

// TestLoadBothConfigAndPath tests error when both --from-config and path provided
func TestLoadBothConfigAndPath(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{"secret/ssh/test"},
	}
	provider := &mockProviderForLoad{}

	err := Load(provider, cfg, []string{"--from-config", "secret/ssh/github"})

	if err == nil {
		t.Error("expected error for both --from-config and path, got nil")
	}
	if err != nil && err.Error() != "cannot use both --from-config and direct path" {
		t.Errorf("expected 'cannot use both' error, got: %v", err)
	}
}

// TestLoadUnknownFlag tests error for unknown flag
func TestLoadUnknownFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{"secret/ssh/test"},
	}
	provider := &mockProviderForLoad{}

	err := Load(provider, cfg, []string{"--unknown-flag"})

	if err == nil {
		t.Error("expected error for unknown flag, got nil")
	}
}
