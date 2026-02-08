package config

import (
	"encoding/json"
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

		paths := cfg.GetPaths()
		if len(paths) != 1 || paths[0] != "secret/ssh/github" {
			t.Errorf("GetPaths() = %v, want [secret/ssh/github]", paths)
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

func TestGetPaths(t *testing.T) {
	cfg := &Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{"path1", "path2"},
	}

	paths := cfg.GetPaths()
	if len(paths) != 2 {
		t.Errorf("GetPaths() returned %d paths, want 2", len(paths))
	}

	cfg2 := &Config{VaultPaths: nil}
	paths2 := cfg2.GetPaths()
	if len(paths2) != 0 {
		t.Errorf("GetPaths() with nil returned %d paths, want 0", len(paths2))
	}
}

func TestConfigAddPath(t *testing.T) {
	cfg := &Config{
		DefaultProvider: ProviderVault,
		VaultPaths:      []string{"secret/ssh/existing"},
	}

	// Test adding new path
	err := cfg.AddPath("secret/ssh/new")
	if err != nil {
		t.Errorf("AddPath failed: %v", err)
	}

	if len(cfg.VaultPaths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(cfg.VaultPaths))
	}

	if cfg.VaultPaths[1] != "secret/ssh/new" {
		t.Errorf("Expected new path, got %q", cfg.VaultPaths[1])
	}
}

func TestConfigAddPath_duplicate(t *testing.T) {
	// Use temp file for testing
	tmpDir := t.TempDir()
	tmpConfigDir := filepath.Join(tmpDir, ".config")
	os.Mkdir(tmpConfigDir, 0755)
	tmpConfigPath := filepath.Join(tmpConfigDir, "sm-ssh-add.json")

	cfg := &Config{
		DefaultProvider: ProviderVault,
		VaultPaths:      []string{"secret/ssh/existing"},
	}

	// Set up temp home directory
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Test adding duplicate path - should return nil (no-op)
	err := cfg.AddPath("secret/ssh/existing")
	if err != nil {
		t.Errorf("Expected no error for duplicate path, got: %v", err)
	}

	// Verify no duplicate added
	if len(cfg.VaultPaths) != 1 {
		t.Errorf("Expected 1 path, got %d", len(cfg.VaultPaths))
	}

	// Test adding new path writes to file
	err = cfg.AddPath("secret/ssh/new")
	if err != nil {
		t.Errorf("AddPath failed: %v", err)
	}

	if len(cfg.VaultPaths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(cfg.VaultPaths))
	}

	// Verify file was written
	if _, err := os.Stat(tmpConfigPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(tmpConfigPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var readCfg Config
	if err := json.Unmarshal(data, &readCfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if len(readCfg.VaultPaths) != 2 {
		t.Errorf("Expected 2 paths in file, got %d", len(readCfg.VaultPaths))
	}
}
