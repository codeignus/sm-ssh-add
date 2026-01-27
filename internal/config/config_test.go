package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	// Set up test config directory
	origHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	testConfigDir := filepath.Join(tempDir, ".config")
	os.Mkdir(testConfigDir, 0755)
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", origHome)

	configPath := filepath.Join(testConfigDir, ConfigFileName)

	t.Run("valid config", func(t *testing.T) {
		os.WriteFile(configPath, []byte(`{"default_provider": "vault", "vault_paths": ["secret/ssh/github"]}`), 0644)

		cfg, err := Read()
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}

		if cfg.DefaultProvider != "vault" {
			t.Errorf("DefaultProvider = %q, want vault", cfg.DefaultProvider)
		}

		paths := cfg.GetVaultPaths()
		if len(paths) != 1 || paths[0] != "secret/ssh/github" {
			t.Errorf("GetVaultPaths() = %v, want [secret/ssh/github]", paths)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		os.WriteFile(configPath, []byte(`{invalid`), 0644)

		_, err := Read()
		if err == nil {
			t.Fatal("Read() should return error for invalid JSON")
		}
	})

	t.Run("empty default_provider", func(t *testing.T) {
		os.WriteFile(configPath, []byte(`{"default_provider": ""}`), 0644)

		_, err := Read()
		if err != ErrEmptyProvider {
			t.Errorf("Read() error = %v, want ErrEmptyProvider", err)
		}
	})

	t.Run("invalid provider", func(t *testing.T) {
		os.WriteFile(configPath, []byte(`{"default_provider": "foobar"}`), 0644)

		_, err := Read()
		if err != ErrInvalidProvider {
			t.Errorf("Read() error = %v, want ErrInvalidProvider", err)
		}
	})
}

func TestGetVaultPaths(t *testing.T) {
	cfg := &Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{"path1", "path2"},
	}

	paths := cfg.GetVaultPaths()
	if len(paths) != 2 {
		t.Errorf("GetVaultPaths() returned %d paths, want 2", len(paths))
	}

	cfg2 := &Config{VaultPaths: nil}
	paths2 := cfg2.GetVaultPaths()
	if len(paths2) != 0 {
		t.Errorf("GetVaultPaths() with nil returned %d paths, want 0", len(paths2))
	}
}
