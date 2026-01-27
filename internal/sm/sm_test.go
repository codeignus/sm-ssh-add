package sm

import (
	"testing"
)

func TestKeyValue_Struct(t *testing.T) {
	tests := []struct {
		name            string
		kv              KeyValue
		wantPrivateKey  bool
		wantPublicKey   bool
		wantRequirePass bool
	}{
		{
			name: "key with passphrase",
			kv: KeyValue{
				PrivateKey:        []byte("private-key-data"),
				PublicKey:         []byte("public-key-data"),
				RequirePassphrase: true,
			},
			wantPrivateKey:  true,
			wantPublicKey:   true,
			wantRequirePass: true,
		},
		{
			name: "key without passphrase",
			kv: KeyValue{
				PrivateKey:        []byte("private-key-data"),
				PublicKey:         []byte("public-key-data"),
				RequirePassphrase: false,
			},
			wantPrivateKey:  true,
			wantPublicKey:   true,
			wantRequirePass: false,
		},
		{
			name: "empty keys",
			kv: KeyValue{
				PrivateKey:        []byte{},
				PublicKey:         []byte{},
				RequirePassphrase: false,
			},
			wantPrivateKey:  false,
			wantPublicKey:   false,
			wantRequirePass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasPrivateKey := len(tt.kv.PrivateKey) > 0
			hasPublicKey := len(tt.kv.PublicKey) > 0

			if hasPrivateKey != tt.wantPrivateKey {
				t.Errorf("PrivateKey presence = %v, want %v", hasPrivateKey, tt.wantPrivateKey)
			}
			if hasPublicKey != tt.wantPublicKey {
				t.Errorf("PublicKey presence = %v, want %v", hasPublicKey, tt.wantPublicKey)
			}
			if tt.kv.RequirePassphrase != tt.wantRequirePass {
				t.Errorf("RequirePassphrase = %v, want %v", tt.kv.RequirePassphrase, tt.wantRequirePass)
			}
		})
	}
}

func TestGet_UnsupportedProvider(t *testing.T) {
	_, err := Get("aws", "some/path")
	if err == nil {
		t.Error("expected error for unsupported provider, got nil")
	}
}

func TestStore_UnsupportedProvider(t *testing.T) {
	kv := &KeyValue{
		PrivateKey:        []byte("test"),
		PublicKey:         []byte("test"),
		RequirePassphrase: false,
	}
	err := Store("aws", "some/path", kv)
	if err == nil {
		t.Error("expected error for unsupported provider, got nil")
	}
}

func TestProviderConstants(t *testing.T) {
	if ProviderVault != "vault" {
		t.Errorf("ProviderVault = %s, want 'vault'", ProviderVault)
	}
}
