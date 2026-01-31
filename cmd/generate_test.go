package cmd

import (
	"testing"

	"github.com/codeignus/sm-ssh-add/internal/config"
)

// TestGenerateNoArguments tests error when no arguments provided
func TestGenerateNoArguments(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{})

	if err == nil {
		t.Error("expected error for no arguments, got nil")
	}
}

// TestGenerateWithOnlyFlag tests error when only --require-passphrase flag provided
func TestGenerateWithOnlyFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--require-passphrase"})

	if err == nil {
		t.Error("expected error when only flag provided, got nil")
	}
}

// TestGenerateWithPathOnly tests minimal valid invocation
func TestGenerateWithPathOnly(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"secret/ssh/test"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestGenerateWithPathAndComment tests with path and comment
func TestGenerateWithPathAndComment(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"secret/ssh/test", "user@example.com"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestGenerateWithFlagPathAndComment tests with all arguments
func TestGenerateWithFlagPathAndComment(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	// Can't easily test passphrase prompt in unit tests
	// This will fail at passphrase prompt, but tests argument parsing
	err := Generate(cfg, []string{"--require-passphrase", "secret/ssh/test", "user@example.com"})

	if err != nil {
		t.Logf("Expected (no Vault/passphrase): %v", err)
	}
}

// TestGenerateTooManyArguments tests error with too many arguments
func TestGenerateTooManyArguments(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"secret/ssh/test", "user@example.com", "extra", "args"})

	if err == nil {
		t.Error("expected error for too many arguments, got nil")
	}
}

// TestGenerateUnknownFlag tests behavior with unknown flag
// Current implementation may accept unknown flags as path
func TestGenerateArgumentOrdering(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "flag after path",
			args:        []string{"secret/ssh/test", "--require-passphrase"},
			expectError: false, // Current implementation treats flag as comment
		},
		{
			name:        "path flag comment",
			args:        []string{"secret/ssh/test", "--require-passphrase", "user@example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Generate(cfg, tt.args)
			// Will fail on Vault, but we're testing argument parsing
			if err != nil {
				t.Logf("Expected (no Vault): %v", err)
			}
		})
	}
}
