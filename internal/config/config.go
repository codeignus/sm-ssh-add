package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/codeignus/sm-ssh-add/internal/sm"
)

// ConfigFileName is the name of the config file
const ConfigFileName = "sm-ssh-add.json"

// Config holds the application configuration. It reads from ~/.config/sm-ssh-add.json
// and contains the default provider and secret manager paths to load keys from.
type Config struct {
	DefaultProvider string   `json:"default_provider"`
	VaultPaths      []string `json:"vault_paths,omitempty"`
}

// Read reads and parses the config file from ~/.config/sm-ssh-add.json
func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, wrapError(err, "failed to get home directory")
	}

	configDir := filepath.Join(homeDir, ".config")
	configPath := filepath.Join(configDir, ConfigFileName)

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
	case sm.ProviderVault:
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
