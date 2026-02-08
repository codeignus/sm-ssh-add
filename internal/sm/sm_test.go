package sm

import (
	"errors"
	"testing"

	"github.com/codeignus/sm-ssh-add/internal/config"
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

func TestInitProvider_UnsupportedProvider(t *testing.T) {
	cfg := &config.Config{DefaultProvider: "aws"}
	_, err := InitProvider(cfg)
	if err == nil {
		t.Error("expected error for unsupported provider, got nil")
	}
}

func TestProviderConstants(t *testing.T) {
	if config.ProviderVault != "vault" {
		t.Errorf("ProviderVault = %s, want 'vault'", config.ProviderVault)
	}
}

// mockProvider is a mock implementation of Provider interface for testing
type mockProvider struct {
	getFunc   func(path string) (*KeyValue, error)
	storeFunc func(path string, kv *KeyValue) error
	checkFunc func(path string) (bool, error)
}

func (m *mockProvider) Get(path string) (*KeyValue, error) {
	if m.getFunc != nil {
		return m.getFunc(path)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProvider) Store(path string, kv *KeyValue) error {
	if m.storeFunc != nil {
		return m.storeFunc(path, kv)
	}
	return errors.New("not implemented")
}

func (m *mockProvider) CheckExists(path string) (bool, error) {
	if m.checkFunc != nil {
		return m.checkFunc(path)
	}
	return false, errors.New("not implemented")
}

func TestProviderInterface_Get(t *testing.T) {
	tests := []struct {
		name    string
		getFunc func(path string) (*KeyValue, error)
		wantErr bool
	}{
		{
			name: "successful get",
			getFunc: func(path string) (*KeyValue, error) {
				return &KeyValue{
					PrivateKey:        []byte("test-private"),
					PublicKey:         []byte("test-public"),
					RequirePassphrase: false,
				}, nil
			},
			wantErr: false,
		},
		{
			name: "get returns error",
			getFunc: func(path string) (*KeyValue, error) {
				return nil, errors.New("not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{getFunc: tt.getFunc}
			kv, err := mock.Get("test/path")

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantErr && kv != nil {
				if len(kv.PrivateKey) == 0 {
					t.Error("expected PrivateKey to be set")
				}
				if len(kv.PublicKey) == 0 {
					t.Error("expected PublicKey to be set")
				}
			}
		})
	}
}

func TestProviderInterface_Store(t *testing.T) {
	tests := []struct {
		name      string
		storeFunc func(path string, kv *KeyValue) error
		wantErr   bool
	}{
		{
			name: "successful store",
			storeFunc: func(path string, kv *KeyValue) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "store returns error",
			storeFunc: func(path string, kv *KeyValue) error {
				return errors.New("store failed")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{storeFunc: tt.storeFunc}
			kv := &KeyValue{
				PrivateKey:        []byte("test-private"),
				PublicKey:         []byte("test-public"),
				RequirePassphrase: false,
			}
			err := mock.Store("test/path", kv)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestProviderInterface_CheckExists(t *testing.T) {
	tests := []struct {
		name       string
		checkFunc  func(path string) (bool, error)
		wantExists bool
		wantErr    bool
	}{
		{
			name: "path exists",
			checkFunc: func(path string) (bool, error) {
				return true, nil
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name: "path does not exist",
			checkFunc: func(path string) (bool, error) {
				return false, nil
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "check returns error",
			checkFunc: func(path string) (bool, error) {
				return false, errors.New("check failed")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{checkFunc: tt.checkFunc}
			exists, err := mock.CheckExists("test/path")

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantErr && exists != tt.wantExists {
				t.Errorf("exists = %v, want %v", exists, tt.wantExists)
			}
		})
	}
}

// Test that VaultClient correctly implements Provider interface
// This uses the real VaultClient struct but verifies it can be used as Provider
func TestVaultClient_ImplementsProvider(t *testing.T) {
	// This is a compile-time check that VaultClient implements Provider
	// If VaultClient doesn't implement Provider, this won't compile
	var _ Provider = (*VaultClient)(nil)

	// Note: The var assignment above is a compile-time check only.
	// Actual connection tests require Vault running and are in integration tests.
}

// mockVaultClient is a minimal test double that mimics VaultClient behavior
// without requiring actual Vault connection
type mockVaultClient struct {
	getCalled   bool
	storeCalled bool
	checkCalled bool
	lastPath    string
	lastKV      *KeyValue
	returnError bool
}

func (m *mockVaultClient) Get(path string) (*KeyValue, error) {
	m.getCalled = true
	m.lastPath = path
	if m.returnError {
		return nil, ErrPathNotFound
	}
	return &KeyValue{
		PrivateKey:        []byte("mock-private-key"),
		PublicKey:         []byte("mock-public-key"),
		RequirePassphrase: false,
	}, nil
}

func (m *mockVaultClient) Store(path string, kv *KeyValue) error {
	m.storeCalled = true
	m.lastPath = path
	m.lastKV = kv
	if m.returnError {
		return errors.New("mock store error")
	}
	return nil
}

func (m *mockVaultClient) CheckExists(path string) (bool, error) {
	m.checkCalled = true
	m.lastPath = path
	if m.returnError {
		return false, errors.New("mock check error")
	}
	return true, nil
}

// TestVaultClient_VerifyMethodCalls verifies that Provider interface methods
// are properly called on VaultClient-like implementation
func TestVaultClient_VerifyMethodCalls(t *testing.T) {
	mock := &mockVaultClient{}

	// Test Get method through Provider interface
	var provider Provider = mock
	kv, err := provider.Get("secret/ssh/key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !mock.getCalled {
		t.Error("Get was not called on underlying client")
	}
	if mock.lastPath != "secret/ssh/key" {
		t.Errorf("Get path = %s, want secret/ssh/key", mock.lastPath)
	}
	if kv == nil {
		t.Error("Get returned nil KeyValue")
	}

	// Test Store method through Provider interface
	mock.getCalled = false // reset
	testKV := &KeyValue{
		PrivateKey:        []byte("test-private"),
		PublicKey:         []byte("test-public"),
		RequirePassphrase: false,
	}
	err = provider.Store("secret/ssh/new", testKV)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}
	if !mock.storeCalled {
		t.Error("Store was not called on underlying client")
	}
	if mock.lastPath != "secret/ssh/new" {
		t.Errorf("Store path = %s, want secret/ssh/new", mock.lastPath)
	}

	// Test CheckExists method through Provider interface
	mock.storeCalled = false // reset
	exists, err := provider.CheckExists("secret/ssh/key")
	if err != nil {
		t.Fatalf("CheckExists failed: %v", err)
	}
	if !mock.checkCalled {
		t.Error("CheckExists was not called on underlying client")
	}
	if !exists {
		t.Error("CheckExists returned false, want true")
	}
}
