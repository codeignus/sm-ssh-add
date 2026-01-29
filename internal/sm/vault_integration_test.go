//go:build integration

package sm

import (
	"testing"
)

func TestStoreKV_stores_key_value_data_successfully(t *testing.T) {
	// Setup: Create VaultClient (will use VAULT_ADDR/VAULT_TOKEN from env)
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	// Create test KeyValue with passphrase
	kv := &KeyValue{
		PrivateKey:        []byte("test-private-key"),
		PublicKey:         []byte("test-public-key"),
		RequirePassphrase: true,
	}

	// Test: Store the data
	testPath := "secret/ssh/test-store-success"
	err = client.StoreKV(testPath, kv)

	// Verify: Should not return error
	if err != nil {
		t.Errorf("StoreKV failed: %v", err)
	}

	// Cleanup: Remove test data
	client.client.Logical().Delete(testPath)
}

func TestStoreKV_stores_key_without_passphrase(t *testing.T) {
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	kv := &KeyValue{
		PrivateKey:        []byte("test-private-key-no-pass"),
		PublicKey:         []byte("test-public-key-no-pass"),
		RequirePassphrase: false,
	}

	testPath := "secret/ssh/test-store-no-passphrase"
	err = client.StoreKV(testPath, kv)

	if err != nil {
		t.Errorf("StoreKV without passphrase failed: %v", err)
	}

	client.client.Logical().Delete(testPath)
}

func TestStoreKV_rejects_empty_path(t *testing.T) {
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	kv := &KeyValue{
		PrivateKey:        []byte("test-private-key"),
		PublicKey:         []byte("test-public-key"),
		RequirePassphrase: false,
	}

	err = client.StoreKV("", kv)

	if err == nil {
		t.Error("Expected error when storing to empty path, got nil")
	}
}

func TestGetKV_retrieves_stored_key_value_data(t *testing.T) {
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	// Setup: Store known data
	originalKV := &KeyValue{
		PrivateKey:        []byte("test-private-retrieve"),
		PublicKey:         []byte("test-public-retrieve"),
		RequirePassphrase: true,
	}
	testPath := "secret/ssh/test-retrieve"
	err = client.StoreKV(testPath, originalKV)
	if err != nil {
		t.Fatalf("Setup failed: StoreKV error: %v", err)
	}

	// Test: Retrieve the data
	retrievedKV, err := client.GetKV(testPath)
	if err != nil {
		t.Errorf("GetKV failed: %v", err)
	}

	// Verify: All fields match
	if string(retrievedKV.PrivateKey) != string(originalKV.PrivateKey) {
		t.Errorf("PrivateKey mismatch: got %q, want %q", retrievedKV.PrivateKey, originalKV.PrivateKey)
	}
	if string(retrievedKV.PublicKey) != string(originalKV.PublicKey) {
		t.Errorf("PublicKey mismatch: got %q, want %q", retrievedKV.PublicKey, originalKV.PublicKey)
	}
	if retrievedKV.RequirePassphrase != originalKV.RequirePassphrase {
		t.Errorf("RequirePassphrase mismatch: got %v, want %v", retrievedKV.RequirePassphrase, originalKV.RequirePassphrase)
	}

	// Cleanup
	client.client.Logical().Delete(testPath)
}

func TestIntegration_works_with_hashicorp_vault(t *testing.T) {
	// This test verifies full workflow with HashiCorp Vault
	// Environment: VAULT_ADDR and VAULT_TOKEN must be set
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	// Test data
	kv := &KeyValue{
		PrivateKey:        []byte("vault-test-private"),
		PublicKey:         []byte("vault-test-public"),
		RequirePassphrase: true,
	}
	testPath := "secret/ssh/vault-e2e-test"

	// Store
	err = client.StoreKV(testPath, kv)
	if err != nil {
		t.Errorf("StoreKV failed: %v", err)
	}

	// Retrieve
	retrieved, err := client.GetKV(testPath)
	if err != nil {
		t.Errorf("GetKV failed: %v", err)
	}

	// Verify round-trip
	if string(retrieved.PrivateKey) != string(kv.PrivateKey) {
		t.Error("PrivateKey round-trip failed")
	}
	if string(retrieved.PublicKey) != string(kv.PublicKey) {
		t.Error("PublicKey round-trip failed")
	}
	if retrieved.RequirePassphrase != kv.RequirePassphrase {
		t.Error("RequirePassphrase round-trip failed")
	}

	// Cleanup
	client.client.Logical().Delete(testPath)
}

func TestIntegration_works_with_openbao(t *testing.T) {
	// This test verifies full workflow with OpenBao
	// Environment: BAO_ADDR and BAO_TOKEN must be set
	// Note: This will only run in the OpenBao job of the integration workflow
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create OpenBao client: %v", err)
	}

	// Test data
	kv := &KeyValue{
		PrivateKey:        []byte("bao-test-private"),
		PublicKey:         []byte("bao-test-public"),
		RequirePassphrase: false,
	}
	testPath := "secret/ssh/bao-e2e-test"

	// Store
	err = client.StoreKV(testPath, kv)
	if err != nil {
		t.Errorf("StoreKV failed: %v", err)
	}

	// Retrieve
	retrieved, err := client.GetKV(testPath)
	if err != nil {
		t.Errorf("GetKV failed: %v", err)
	}

	// Verify round-trip
	if string(retrieved.PrivateKey) != string(kv.PrivateKey) {
		t.Error("PrivateKey round-trip failed")
	}
	if string(retrieved.PublicKey) != string(kv.PublicKey) {
		t.Error("PublicKey round-trip failed")
	}
	if retrieved.RequirePassphrase != kv.RequirePassphrase {
		t.Error("RequirePassphrase round-trip failed")
	}

	// Cleanup
	client.client.Logical().Delete(testPath)
}

func TestGetKV_returns_path_not_found_error_for_nonexistent_path(t *testing.T) {
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	// Test: Try to get from path that doesn't exist
	_, err = client.GetKV("secret/ssh/does-not-exist-12345")

	// Verify: Should return ErrPathNotFound
	if err != ErrPathNotFound {
		t.Errorf("Expected ErrPathNotFound, got: %v", err)
	}
}

func TestGetKV_rejects_malformed_vault_data_missing_private_key(t *testing.T) {
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	// Setup: Manually write malformed data (missing private_key)
	testPath := "secret/ssh/test-malformed-no-private"
	malformedData := map[string]interface{}{
		"data": map[string]interface{}{
			"public_key":         "test-public",
			"require_passphrase": false,
		},
	}
	_, err = client.client.Logical().Write(testPath, malformedData)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Test: Try to get the malformed data
	_, err = client.GetKV(testPath)

	// Verify: Should return ErrInvalidKeyFormat
	if err != ErrInvalidKeyFormat {
		t.Errorf("Expected ErrInvalidKeyFormat, got: %v", err)
	}

	// Cleanup
	client.client.Logical().Delete(testPath)
}

func TestGetKV_rejects_malformed_vault_data_missing_public_key(t *testing.T) {
	client, err := NewVaultClient()
	if err != nil {
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	// Setup: Manually write malformed data (missing public_key)
	testPath := "secret/ssh/test-malformed-no-public"
	malformedData := map[string]interface{}{
		"data": map[string]interface{}{
			"private_key":        "test-private",
			"require_passphrase": false,
		},
	}
	_, err = client.client.Logical().Write(testPath, malformedData)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Test: Try to get the malformed data
	_, err = client.GetKV(testPath)

	// Verify: Should return ErrInvalidKeyFormat
	if err != ErrInvalidKeyFormat {
		t.Errorf("Expected ErrInvalidKeyFormat, got: %v", err)
	}

	// Cleanup
	client.client.Logical().Delete(testPath)
}
