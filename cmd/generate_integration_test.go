//go:build integration

package cmd

import (
	"os"
	"slices"
	"testing"

	"github.com/codeignus/sm-ssh-add/internal/config"
	"github.com/codeignus/sm-ssh-add/internal/sm"
)

func TestGenerateCommand_RegenerateFlag_OverwritesExistingKey(t *testing.T) {
	cfg := &config.Config{DefaultProvider: config.ProviderVault}

	// Initialize provider
	provider, err := sm.InitProvider(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize provider: %v", err)
	}

	// Setup: Store initial key
	path := "secret/data/ssh/test-regenerate"
	initialKV := &sm.KeyValue{
		PrivateKey:        []byte("initial-private-key"),
		PublicKey:         []byte("initial-public-key"),
		RequirePassphrase: false,
		Comment:           "initial@test",
	}
	err = provider.Store(path, initialKV)
	if err != nil {
		t.Fatalf("Setup failed: Store error: %v", err)
	}

	// Test: Generate with --regenerate
	args := []string{"--regenerate", path, "regenerated@test"}
	err = Generate(provider, cfg, args)

	// Verify: No error
	if err != nil {
		t.Errorf("Generate with --regenerate failed: %v", err)
	}

	// Verify: Key was overwritten
	retrieved, err := provider.Get(path)
	if err != nil {
		t.Fatalf("Failed to retrieve key: %v", err)
	}
	if retrieved.Comment == "initial@test" {
		t.Error("Key was not overwritten - comment still has initial value")
	}
	if retrieved.Comment != "regenerated@test" {
		t.Errorf("Comment not updated correctly: got %q, want %q", retrieved.Comment, "regenerated@test")
	}
}

func TestGenerateCommand_SavePathFlag_AddsPathToConfig(t *testing.T) {
	// Setup: Use temp directory for config
	origHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", origHome)

	// Create initial config file
	configDir := tempDir + "/.config"
	os.Mkdir(configDir, 0755)
	configPath := configDir + "/sm-ssh-add.json"
	initialConfig := `{"default_provider": "vault", "vault_paths": ["secret/ssh/existing"]}`
	os.WriteFile(configPath, []byte(initialConfig), 0600)

	// Reload config to pick up initial paths
	cfg, err := config.Read()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Initialize provider
	provider, err := sm.InitProvider(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize provider: %v", err)
	}

	// Test: Generate with --save-path
	testPath := "secret/data/ssh/test-save-path"
	args := []string{"--save-path", testPath, "savepath@test"}
	err = Generate(provider, cfg, args)

	// Verify: No error
	if err != nil {
		t.Errorf("Generate with --save-path failed: %v", err)
	}

	// Verify: Path was added to in-memory config
	if !slices.Contains(cfg.VaultPaths, testPath) {
		t.Errorf("Path %q not found in config.VaultPaths: %v", testPath, cfg.VaultPaths)
	}

	// Verify: Path was written to config file on disk
	reloadedCfg, err := config.Read()
	if err != nil {
		t.Fatalf("Failed to re-read config: %v", err)
	}
	if !slices.Contains(reloadedCfg.VaultPaths, testPath) {
		t.Errorf("Path %q not found in reloaded config from disk: %v", testPath, reloadedCfg.VaultPaths)
	}
}

func TestGenerateCommand_SavePathFlag_DuplicateIsNoOp(t *testing.T) {
	// Setup: Use temp directory for config
	origHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", origHome)

	// Create config with existing path
	configDir := tempDir + "/.config"
	os.Mkdir(configDir, 0755)
	configPath := configDir + "/sm-ssh-add.json"
	existingPath := "secret/data/ssh/test-duplicate-save"
	initialConfig := `{"default_provider": "vault", "vault_paths": ["` + existingPath + `"]}`
	os.WriteFile(configPath, []byte(initialConfig), 0600)

	cfg, err := config.Read()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Verify initial state
	initialCount := len(cfg.VaultPaths)
	if initialCount != 1 {
		t.Fatalf("Expected 1 path initially, got %d", initialCount)
	}

	// Test: Add same path again with --save-path (simulate duplicate scenario)
	// Since AddPath is a no-op for duplicates, we test it directly
	err = cfg.AddPath(existingPath)
	if err != nil {
		t.Errorf("AddPath with duplicate should not error, got: %v", err)
	}

	// Verify: Path count unchanged
	if len(cfg.VaultPaths) != initialCount {
		t.Errorf("Path count changed after duplicate AddPath: got %d, want %d", len(cfg.VaultPaths), initialCount)
	}
}

func TestGenerateCommand_WithoutRegenerate_RejectsDuplicatePath(t *testing.T) {
	cfg := &config.Config{DefaultProvider: config.ProviderVault}

	// Initialize provider
	provider, err := sm.InitProvider(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize provider: %v", err)
	}

	// Setup: Store a key
	path := "secret/data/ssh/test-no-regenerate"
	initialKV := &sm.KeyValue{
		PrivateKey:        []byte("initial-private-key"),
		PublicKey:         []byte("initial-public-key"),
		RequirePassphrase: false,
		Comment:           "existing@test",
	}
	err = provider.Store(path, initialKV)
	if err != nil {
		t.Fatalf("Setup failed: Store error: %v", err)
	}

	// Test: Generate WITHOUT --regenerate
	args := []string{path, "should-fail@test"}
	err = Generate(provider, cfg, args)

	// Verify: Error about key existing
	if err == nil {
		t.Error("Expected error when generating duplicate key without --regenerate")
	}
	expectedErrMsg := "key already exists at " + path + " (use --regenerate to overwrite)"
	if err != nil && err.Error() != expectedErrMsg {
		t.Errorf("Wrong error message: got %q, want %q", err.Error(), expectedErrMsg)
	}

	// Verify: Original key was not modified
	retrieved, err := provider.Get(path)
	if err != nil {
		t.Fatalf("Failed to retrieve key: %v", err)
	}
	if retrieved.Comment != "existing@test" {
		t.Errorf("Original key was modified: comment = %q, want %q", retrieved.Comment, "existing@test")
	}
}

func TestGenerateCommand_NewPath_Succeeds(t *testing.T) {
	cfg := &config.Config{DefaultProvider: config.ProviderVault}

	// Initialize provider
	provider, err := sm.InitProvider(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize provider: %v", err)
	}

	// Test: Generate to a new path
	path := "secret/data/ssh/test-new-path-generate"
	args := []string{path, "newkey@test"}
	err = Generate(provider, cfg, args)

	// Verify: No error
	if err != nil {
		t.Errorf("Generate to new path failed: %v", err)
	}

	// Verify: Key was stored
	retrieved, err := provider.Get(path)
	if err != nil {
		t.Fatalf("Failed to retrieve key: %v", err)
	}
	if retrieved.Comment != "newkey@test" {
		t.Errorf("Comment incorrect: got %q, want %q", retrieved.Comment, "newkey@test")
	}
}
