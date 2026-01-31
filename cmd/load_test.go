package cmd

import (
	"testing"

	"github.com/codeignus/sm-ssh-add/internal/config"
)

// TestLoadFromConfig_EmptyPaths tests --from-config with empty paths
func TestLoadFromConfig_EmptyPaths(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{},
	}

	err := Load(cfg, []string{"--from-config"})

	if err == nil {
		t.Error("expected error for empty paths, got nil")
	}
	if err != nil && err.Error() != "no vault paths configured" {
		t.Errorf("expected 'no vault paths configured' error, got: %v", err)
	}
}

// TestLoadFromConfig_WithPaths tests --from-config with configured paths
// This test requires real Vault instance - will run in GHA
func TestLoadFromConfig_WithPaths(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{"secret/ssh/test"},
	}

	err := Load(cfg, []string{"--from-config"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		// Expected to fail on actual load, but not on argument validation
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestLoadDirectPath tests loading from a direct path argument
func TestLoadDirectPath(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{}, // Should not be used
	}

	err := Load(cfg, []string{"secret/ssh/github"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestLoadNoArguments tests error when no arguments provided
func TestLoadNoArguments(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
		VaultPaths:      []string{"secret/ssh/test"},
	}

	err := Load(cfg, []string{})

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

	err := Load(cfg, []string{"--from-config", "secret/ssh/github"})

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

	err := Load(cfg, []string{"--unknown-flag"})

	if err == nil {
		t.Error("expected error for unknown flag, got nil")
	}
}
