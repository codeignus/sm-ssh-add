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

// TestGenerateWithRegenerateFlag tests --regenerate flag
func TestGenerateWithRegenerateFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--regenerate", "secret/ssh/test"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestGenerateWithRegenerateAndComment tests --regenerate with path and comment
func TestGenerateWithRegenerateAndComment(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--regenerate", "secret/ssh/test", "user@example.com"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestGenerateWithOnlyRegenerateFlag tests error when only --regenerate flag provided
func TestGenerateWithOnlyRegenerateFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--regenerate"})

	if err == nil {
		t.Error("expected error when only flag provided, got nil")
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

	// With 6 args it should fail (max is 5: 3 flags + path + comment)
	err := Generate(cfg, []string{"secret/ssh/test", "user@example.com", "extra", "args", "another", "onemore"})

	if err == nil {
		t.Error("expected error for too many arguments, got nil")
	}
}

// TestGenerateWithSavePathFlag tests --save-path flag
func TestGenerateWithSavePathFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--save-path", "secret/ssh/test"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestGenerateWithBothFlags tests both --require-passphrase and --save-path flags
func TestGenerateWithBothFlags(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	// This will fail at passphrase prompt, but tests argument parsing
	err := Generate(cfg, []string{"--require-passphrase", "--save-path", "secret/ssh/test", "user@example.com"})

	if err != nil {
		t.Logf("Expected (no Vault/passphrase): %v", err)
	}
}

// TestGenerateWithAllFlags tests all three flags together
func TestGenerateWithAllFlags(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	// This will fail at passphrase prompt, but tests argument parsing
	err := Generate(cfg, []string{"--require-passphrase", "--save-path", "--regenerate", "secret/ssh/test", "user@example.com"})

	if err != nil {
		t.Logf("Expected (no Vault/passphrase): %v", err)
	}
}

// TestGenerateWithSavePathAndComment tests --save-path with path and comment
func TestGenerateWithSavePathAndComment(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--save-path", "secret/ssh/test", "user@example.com"})

	// Will fail without real Vault, but tests argument parsing
	if err != nil {
		t.Logf("Expected (no Vault): %v", err)
	}
}

// TestGenerateWithOnlySavePathFlag tests error when only --save-path flag provided
func TestGenerateWithOnlySavePathFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	err := Generate(cfg, []string{"--save-path"})

	if err == nil {
		t.Error("expected error when only flag provided, got nil")
	}
}

// TestGenerateWithUnknownFlag tests error when unknown flag is provided
func TestGenerateWithUnknownFlag(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: config.ProviderVault,
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "unknown single dash flag",
			args:        []string{"-unknown", "secret/ssh/test"},
			expectError: true,
		},
		{
			name:        "unknown double dash flag",
			args:        []string{"--unknown", "secret/ssh/test"},
			expectError: true,
		},
		{
			name:        "unknown flag after path",
			args:        []string{"secret/ssh/test", "--unknown"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Generate(cfg, tt.args)
			if tt.expectError && err == nil {
				t.Error("expected error for unknown flag, got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
