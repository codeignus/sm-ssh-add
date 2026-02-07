package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// ConfigFileName is the name of the config file
const ConfigFileName = "sm-ssh-add.json"

// Provider constants
const (
	ProviderVault = "vault"
	// ProviderAWS = "aws" // Future implementation
)

// Config holds the application configuration. It reads from ~/.config/sm-ssh-add.json
// and contains the default provider and secret manager paths to load keys from.
type Config struct {
	DefaultProvider    string   `json:"default_provider"`
	VaultPaths         []string `json:"vault_paths,omitempty"`
	VaultApproleRoleID string   `json:"vault_approle_role_id,omitempty"` // If set, use Vault Approle auth instead of token
}

// GetVaultApproleRoleID returns the configured Vault Approle Role ID.
func (c *Config) GetVaultApproleRoleID() string {
	return c.VaultApproleRoleID
}

// getConfigFilePath returns the path to the config file
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", ConfigFileName), nil
}

// Read reads and parses the config file from ~/.config/sm-ssh-add.json
func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, wrapError(err, "failed to get config path")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrConfigFileNotFound
		}
		return nil, wrapError(err, "failed to read config file")
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, wrapError(err, "failed to parse config file (invalid JSON)")
	}

	// Validate configuration
	if cfg.DefaultProvider == "" {
		return nil, ErrEmptyProvider
	}

	switch cfg.DefaultProvider {
	case ProviderVault:
		// Valid provider
	default:
		return nil, ErrInvalidProvider
	}

	return &cfg, nil
}

// GetVaultPaths returns all configured vault paths
func (c *Config) GetVaultPaths() []string {
	if c.VaultPaths == nil {
		return []string{}
	}
	return c.VaultPaths
}

// AddPath adds a new path to the vault paths list and writes the config file
func (c *Config) AddPath(path string) error {
	// If path already exists, do nothing (no-op)
	if slices.Contains(c.VaultPaths, path) {
		return nil
	}

	c.VaultPaths = append(c.VaultPaths, path)

	// Write the updated config to disk
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0600)
}
