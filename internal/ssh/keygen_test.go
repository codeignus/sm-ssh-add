package ssh

import (
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

// TestGenerateKeyPair_NoPassphrase tests key generation without passphrase
func TestGenerateKeyPair_NoPassphrase(t *testing.T) {
	comment := "test@example.com"
	passphrase := []byte{}

	keyPair, err := GenerateKeyPair(comment, passphrase)
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if keyPair == nil {
		t.Fatal("keyPair is nil")
	}

	if len(keyPair.PrivateKey) == 0 {
		t.Error("private key is empty")
	}

	if len(keyPair.PublicKey) == 0 {
		t.Error("public key is empty")
	}

	// Verify Passphrase field is empty (optional field defaults to empty string)
	if keyPair.Passphrase != "" {
		t.Errorf("Passphrase field should be empty for unencrypted key, got: %q", keyPair.Passphrase)
	}

	// Verify private key can be parsed without passphrase
	signer, err := ssh.ParsePrivateKey(keyPair.PrivateKey)
	if err != nil {
		t.Errorf("failed to parse generated private key: %v", err)
	}

	if signer == nil {
		t.Error("signer is nil after parsing private key")
	}

	// Verify public key is in OpenSSH authorized_keys format
	publicKeyStr := string(keyPair.PublicKey)
	if !strings.HasPrefix(publicKeyStr, "ssh-ed25519") {
		t.Errorf("public key doesn't start with ssh-ed25519: %s", publicKeyStr[:20])
	}

	// Verify comment is stored separately
	if keyPair.Comment != comment {
		t.Errorf("comment field doesn't match: got %q, want %q", keyPair.Comment, comment)
	}
}

// TestGenerateKeyPair_WithPassphrase tests key generation with passphrase
func TestGenerateKeyPair_WithPassphrase(t *testing.T) {
	comment := "secure@example.com"
	passphrase := []byte("test-password-123")

	keyPair, err := GenerateKeyPair(comment, passphrase)
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if keyPair == nil {
		t.Fatal("keyPair is nil")
	}

	if len(keyPair.PrivateKey) == 0 {
		t.Error("private key is empty")
	}

	if len(keyPair.PublicKey) == 0 {
		t.Error("public key is empty")
	}

	// Verify Passphrase field is empty (GenerateKeyPair doesn't store it, only encrypts the private key)
	if keyPair.Passphrase != "" {
		t.Errorf("Passphrase field should be empty after GenerateKeyPair, got: %q", keyPair.Passphrase)
	}

	// Verify private key requires passphrase to parse
	_, err = ssh.ParsePrivateKey(keyPair.PrivateKey)
	if err == nil {
		t.Error("private key should require passphrase but didn't")
	}

	// Verify private key can be parsed with correct passphrase
	signer, err := ssh.ParsePrivateKeyWithPassphrase(keyPair.PrivateKey, passphrase)
	if err != nil {
		t.Errorf("failed to parse private key with correct passphrase: %v", err)
	}

	if signer == nil {
		t.Error("signer is nil after parsing private key with passphrase")
	}

	// Verify wrong passphrase fails
	_, err = ssh.ParsePrivateKeyWithPassphrase(keyPair.PrivateKey, []byte("wrong-password"))
	if err == nil {
		t.Error("private key should fail with wrong passphrase")
	}
}

// TestGenerateKeyPair_EmptyComment tests key generation with empty comment
func TestGenerateKeyPair_EmptyComment(t *testing.T) {
	comment := ""
	passphrase := []byte{}

	keyPair, err := GenerateKeyPair(comment, passphrase)
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if keyPair == nil {
		t.Fatal("keyPair is nil")
	}

	// Should still generate valid keys
	signer, err := ssh.ParsePrivateKey(keyPair.PrivateKey)
	if err != nil {
		t.Errorf("failed to parse generated private key: %v", err)
	}

	if signer == nil {
		t.Error("signer is nil after parsing private key")
	}
}

// TestGenerateKeyPair_Uniqueness tests that each generated key is unique
func TestGenerateKeyPair_Uniqueness(t *testing.T) {
	comment := "test@example.com"
	passphrase := []byte{}

	keyPair1, err := GenerateKeyPair(comment, passphrase)
	if err != nil {
		t.Fatalf("first GenerateKeyPair failed: %v", err)
	}

	keyPair2, err := GenerateKeyPair(comment, passphrase)
	if err != nil {
		t.Fatalf("second GenerateKeyPair failed: %v", err)
	}

	// Private keys should be different
	if string(keyPair1.PrivateKey) == string(keyPair2.PrivateKey) {
		t.Error("generated private keys are identical, expected uniqueness")
	}

	// Public keys should be different
	if string(keyPair1.PublicKey) == string(keyPair2.PublicKey) {
		t.Error("generated public keys are identical, expected uniqueness")
	}
}

// TestGenerateKeyPair_PubKeyFormat tests public key format is valid
func TestGenerateKeyPair_PubKeyFormat(t *testing.T) {
	comment := "test-user@host"
	passphrase := []byte{}

	keyPair, err := GenerateKeyPair(comment, passphrase)
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	// Parse the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(keyPair.PublicKey)
	if err != nil {
		t.Errorf("failed to parse authorized key format: %v", err)
	}

	if pubKey == nil {
		t.Error("parsed public key is nil")
	}

	// Verify it's an ed25519 key
	if pubKey.Type() != "ssh-ed25519" {
		t.Errorf("expected key type ssh-ed25519, got %s", pubKey.Type())
	}
}
