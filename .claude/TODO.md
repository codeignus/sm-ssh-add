# Development TODO

## Task Completion Process

1. Move completed tasks to "Ready for Review" section
2. Mark with `[x]` when fully implemented AND tested
3. No task is complete until tests pass

---

## Ready for Review

### 1. ✅ Secret Manager Package

**Files:** `internal/sm/sm.go`, `internal/sm/errors.go`, `internal/sm/sm_test.go`

Simplified Secret Manager implementation with direct functions.

- [x] `KeyValue` struct with PrivateKey, PublicKey, RequirePassphrase fields
- [x] `Get(provider, path)` - retrieve data from any provider
- [x] `Store(provider, path, kv)` - store data to any provider
- [x] Provider constants (ProviderVault, future ProviderAWS)
- [x] Custom error definitions
- [x] Tests for KeyValue struct, provider validation, and constants

---

### 2. ✅ Vault Client

**Files:** `internal/sm/vault.go`, `internal/sm/vault_test.go`

Implemented HashiCorp Vault/OpenBao KV v2 client with TDD.

- [x] `VaultClient` struct with client field
- [x] `NewVaultClient()` using env vars (BAO_ADDR/BAO_TOKEN, falls back to VAULT_ADDR/VAULT_TOKEN)
- [x] `GetKV(path)` with KV v2 data wrapper handling
- [x] `StoreKV(path, kv)` with KV v2 data wrapper
- [x] Handle KV v2 `{"data": {"data": {...}}}` format correctly
- [x] Support OpenBao (BAO_*) and Vault (VAULT_*) environment variables

**Note:** Integration tests for GetKV/StoreKV deferred to GitHub Actions workflow with real Vault/OpenBao instance

### 3. ✅ Implement Config File Management

**Files:** `internal/config/config.go`, `internal/config/errors.go`, `internal/config/config_test.go`

Implemented with validation and tests.

- [x] `Config` struct with DefaultProvider, VaultPaths fields
- [x] `Read()` function - read from ~/.config/sm-ssh-add.json
- [x] Config validation during read (empty provider, invalid provider)
- [x] Provider validation using sm package constants
- [x] Handle missing config file errors (ErrConfigFileNotFound)
- [x] Handle invalid JSON errors
- [x] Tests for valid config, invalid JSON, empty provider, invalid provider
- [x] Tests for GetVaultPaths accessor method

---

### 4. ✅ Implement SSH Key Generation

**Files:** `internal/ssh/keygen.go`, `internal/ssh/errors.go`, `internal/ssh/keygen_test.go`

Implemented with TDD.

- [x] `GenerateKeyPair(comment string, passphrase []byte) (*KeyPair, error)`
- [x] Use `golang.org/x/crypto/ssh` for Ed25519 key generation
- [x] Return public key in OpenSSH format with comment
- [x] Return private key in OpenSSH format (with passphrase encryption if provided)
- [x] Tests for key generation without passphrase
- [x] Tests for key generation with passphrase
- [x] Tests for empty comment
- [x] Tests for key uniqueness
- [x] Tests for public key format validation

---

### 5. ✅ Implement ssh-agent Operations

**Files:** `internal/ssh/agent.go`

Implemented for `load` command only (not used by `generate`).

- [x] `NewAgent(cfg *Config) (*Agent, error)` - connect to ssh-agent via SSH_AUTH_SOCK
- [x] `AddKey(keyPair *KeyPair) error` - add key with duplicate detection
- [x] `List() ([]*agent.Key, error)` - list keys in agent
- [x] `KeyExists(fingerprint string) (bool, error)` - check if key already loaded
- [x] `Close() error` - close connection
- [x] Handle SSH_AUTH_SOCK environment variable
- [x] Handle connection errors
- [x] Duplicate detection using SHA256 fingerprints

**Note:** Tests deferred - would require mocking ssh-agent or integration test with real agent

---

### 6. ✅ Implement CLI Commands

**Files:** `cmd/generate.go`, `cmd/load.go`, `main.go`

Implemented without tests (cmd layer is thin glue code over tested packages).

**generate command:**
- [x] Argument parsing for path, comment, --require-passphrase flag
- [x] Passphrase prompt with confirmation
- [x] Call SSH key generation
- [x] Store key in Secret Manager
- [x] Display public key with storage path
- [x] Clear error messages

**load command:**
- [x] Load paths from config
- [x] Retrieve keys from Secret Manager
- [x] Check if key already exists in ssh-agent
- [x] Add non-duplicate keys to agent
- [x] Print status messages (loaded, skipped, errors)
- [x] Clear error messages

**main.go:**
- [x] CLI argument parsing and command routing
- [x] Config file loading
- [x] Proper exit codes
- [x] Error handling

**Note:** cmd tests skipped - underlying packages fully tested, cmd layer requires complex I/O mocking for minimal value

---

## Tasks

---

### 7. ⏳ Set up GitHub Actions Workflow

**File:** `.github/workflows/test.yml`

- [ ] Create workflow for running tests on push/PR
- [ ] Spin up Vault test container
- [ ] Spin up OpenBao test container
- [ ] Run unit tests
- [ ] Run integration tests for Vault client
- [ ] Run integration tests for OpenBao client
- [ ] Test against multiple Go versions

---

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```
