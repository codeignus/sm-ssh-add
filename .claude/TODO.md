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

## Tasks

---

### 4. ⏳ Implement SSH Key Generation

**Files:** `internal/ssh/generate.go`, `internal/ssh/generate_test.go`

Follow TDD workflow using `superpowers:test-driven-development` skill.

- [ ] Implement `GenerateKey(comment string, requirePassphrase bool) (*KeyPair, error)`
- [ ] Use `golang.org/x/crypto/ssh` for Ed25519 key generation
- [ ] Handle optional passphrase prompt using terminal stdin
- [ ] Return public key in OpenSSH format
- [ ] Return private key in OpenSSH format (with passphrase if provided)
- [ ] Add tests for key generation without passphrase
- [ ] Add tests for key generation with passphrase
- [ ] Add tests for invalid comment/error handling

---

### 5. ⏳ Implement ssh-agent Operations

**Files:** `internal/ssh/agent.go`, `internal/ssh/agent_test.go`

Follow TDD workflow using `superpowers:test-driven-development` skill.

- [ ] Implement `AddToAgent(key *KeyPair) error` - connect to ssh-agent and add key
- [ ] Implement `ListKeys() ([]KeyInfo, error)` - list keys in agent
- [ ] Implement `KeyExists(fingerprint string) (bool, error)` - check if key already loaded
- [ ] Implement `GetFingerprint(publicKey string) string` - calculate SHA256 fingerprint
- [ ] Handle SSH_AUTH_SOCK environment variable
- [ ] Handle connection errors
- [ ] Add tests for agent operations (may need mock agent interface)

---

### 6. ⏳ Implement CLI Commands

**Files:** `cmd/generate.go`, `cmd/load.go`, `cmd/generate_test.go`, `cmd/load_test.go`

Follow TDD workflow using `superpowers:test-driven-development` skill.

**generate command:**
- [ ] Implement `Generate(path, comment string, requirePassphrase bool, sm SecretManager) error`
- [ ] Call SSH key generation
- [ ] Store key in Secret Manager
- [ ] Print success message with public key
- [ ] Handle errors with clear messages
- [ ] Add tests for successful generation
- [ ] Add tests for various error scenarios

**load command:**
- [ ] Implement `Load(paths []string, sm SecretManager) error`
- [ ] Load each path from Secret Manager
- [ ] Check if key already exists in ssh-agent
- [ ] Add non-duplicate keys to agent
- [ ] Print summary (loaded, skipped, error counts)
- [ ] Handle errors with clear messages
- [ ] Add tests for successful load
- [ ] Add tests for duplicate detection
- [ ] Add tests for various error scenarios

---

### 7. ⏳ Implement CLI Entry Point

**Files:** `internal/cli/cli.go`, `internal/cli/cli_test.go`, `main.go`

Follow TDD workflow using `superpowers:test-driven-development` skill.

- [ ] Implement `Run() error` - parse CLI arguments and route to commands
- [ ] Implement `runGenerate()` - parse generate flags and call cmd.Generate()
- [ ] Implement `runLoad()` - parse load flags and call cmd.Load()
- [ ] Implement `NewSecretManager(provider string) (SecretManager, error)` - factory function
- [ ] Implement `loadConfig() (*Config, error)` - read config file
- [ ] Handle help/version flags
- [ ] Handle invalid commands/arguments
- [ ] Implement proper exit codes
- [ ] Add tests for command routing
- [ ] Add tests for flag parsing
- [ ] Add tests for error handling

---

### 8. ⏳ Set up GitHub Actions Workflow

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
