//go:build integration

package ssh

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/codeignus/sm-ssh-add/internal/config"
)

func TestNewAgent_connects_to_ssh_agent_successfully(t *testing.T) {
	// Setup: Ensure SSH_AUTH_SOCK is set (CI workflow provides this)
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		t.Skip("SSH_AUTH_SOCK not set - skipping agent integration test")
	}

	// Test: Create agent connection
	cfg := &config.Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{},
	}
	agent, err := NewAgent(cfg)

	// Verify: Connection succeeds
	if err != nil {
		t.Fatalf("Failed to connect to ssh-agent: %v", err)
	}
	if agent == nil {
		t.Fatal("Agent is nil but no error returned")
	}

	// Cleanup
	agent.Close()
}

func TestAddKey_adds_key_to_agent_successfully(t *testing.T) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		t.Skip("SSH_AUTH_SOCK not set - skipping agent integration test")
	}

	cfg := &config.Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{},
	}
	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to ssh-agent: %v", err)
	}
	defer agent.Close()

	// Setup: Generate test key pair
	keyPair, err := GenerateKeyPair("test@integration", nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Calculate expected fingerprint
	pubKeyFields := strings.Fields(string(keyPair.PublicKey))
	if len(pubKeyFields) < 2 {
		t.Fatalf("Invalid public key format: %s", string(keyPair.PublicKey))
	}
	keyBytes := []byte(pubKeyFields[1])
	fingerprint := sha256.Sum256(keyBytes)
	expectedFingerprint := base64.RawStdEncoding.EncodeToString(fingerprint[:])

	// Test: Add key to agent
	err = agent.AddKey(keyPair)

	// Verify: No error
	if err != nil {
		t.Errorf("AddKey failed: %v", err)
	}

	// Verify: Key exists in agent
	exists, err := agent.KeyExists(expectedFingerprint)
	if err != nil {
		t.Errorf("KeyExists check failed: %v", err)
	}
	if !exists {
		t.Error("Key was not found in agent after AddKey")
	}

	// Cleanup: Remove key from agent
	agent.client.Remove(keyPair.Signer, keyPair.PublicKey)
}

func TestAddKey_detects_duplicate_key_and_skips(t *testing.T) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		t.Skip("SSH_AUTH_SOCK not set - skipping agent integration test")
	}

	cfg := &config.Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{},
	}
	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to ssh-agent: %v", err)
	}
	defer agent.Close()

	// Setup: Generate and add key first time
	keyPair, err := GenerateKeyPair("duplicate-test@integration", nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	err = agent.AddKey(keyPair)
	if err != nil {
		t.Fatalf("Failed to add key first time: %v", err)
	}

	// Test: Try to add same key again (should skip duplicate)
	err = agent.AddKey(keyPair)

	// Verify: Should succeed (idempotent or skip)
	if err != nil {
		t.Errorf("AddKey duplicate failed: %v", err)
	}

	// Verify: Still only one key in agent (no duplicates)
	keys, err := agent.List()
	if err != nil {
		t.Errorf("Failed to list keys: %v", err)
	}

	// Count matching keys
	pubKeyFields := strings.Fields(string(keyPair.PublicKey))
	keyBytes := []byte(pubKeyFields[1])
	fingerprint := sha256.Sum256(keyBytes)
	expectedFingerprint := base64.RawStdEncoding.EncodeToString(fingerprint[:])

	matchingKeys := 0
	for _, key := range keys {
		if key.String() == expectedFingerprint {
			matchingKeys++
		}
	}

	if matchingKeys != 1 {
		t.Errorf("Expected 1 key in agent, found %d (duplicate detection failed)", matchingKeys)
	}

	// Cleanup
	agent.client.Remove(keyPair.Signer, keyPair.PublicKey)
}

func TestList_lists_keys_from_agent(t *testing.T) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		t.Skip("SSH_AUTH_SOCK not set - skipping agent integration test")
	}

	cfg := &config.Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{},
	}
	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to ssh-agent: %v", err)
	}
	defer agent.Close()

	// Setup: Generate and add known key
	keyPair, err := GenerateKeyPair("list-test@integration", nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	err = agent.AddKey(keyPair)
	if err != nil {
		t.Fatalf("Failed to add key: %v", err)
	}

	// Calculate expected fingerprint
	pubKeyFields := strings.Fields(string(keyPair.PublicKey))
	keyBytes := []byte(pubKeyFields[1])
	fingerprint := sha256.Sum256(keyBytes)
	expectedFingerprint := base64.RawStdEncoding.EncodeToString(fingerprint[:])

	// Test: List keys from agent
	keys, err := agent.List()

	// Verify: No error
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	// Verify: Our key is in the list
	found := false
	for _, key := range keys {
		if key.String() == expectedFingerprint {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Our key not found in agent listing. Expected fingerprint: %s", expectedFingerprint)
	}

	// Cleanup
	agent.client.Remove(keyPair.Signer, keyPair.PublicKey)
}

func TestClose_closes_connection_cleanly(t *testing.T) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		t.Skip("SSH_AUTH_SOCK not set - skipping agent integration test")
	}

	cfg := &config.Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{},
	}
	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to ssh-agent: %v", err)
	}

	// Test: Close connection
	err = agent.Close()

	// Verify: No error on close
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Verify: Client is nil after close
	if agent.client != nil {
		t.Error("Agent client should be nil after Close")
	}
}

func TestEndToEnd_generate_and_load_workflow(t *testing.T) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		t.Skip("SSH_AUTH_SOCK not set - skipping agent integration test")
	}

	cfg := &config.Config{
		DefaultProvider: "vault",
		VaultPaths:      []string{},
	}
	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to ssh-agent: %v", err)
	}
	defer agent.Close()

	// Test: Complete workflow - generate and load
	comment := "e2e-test@integration"
	keyPair, err := GenerateKeyPair(comment, nil)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	err = agent.AddKey(keyPair)
	if err != nil {
		t.Errorf("Failed to load key into agent: %v", err)
	}

	// Verify: Key is actually loaded
	pubKeyFields := strings.Fields(string(keyPair.PublicKey))
	keyBytes := []byte(pubKeyFields[1])
	fingerprint := sha256.Sum256(keyBytes)
	expectedFingerprint := base64.RawStdEncoding.EncodeToString(fingerprint[:])

	exists, err := agent.KeyExists(expectedFingerprint)
	if err != nil {
		t.Errorf("Failed to verify key in agent: %v", err)
	}
	if !exists {
		t.Error("Key was not loaded into agent successfully")
	}

	// Cleanup
	agent.client.Remove(keyPair.Signer, keyPair.PublicKey)
}
