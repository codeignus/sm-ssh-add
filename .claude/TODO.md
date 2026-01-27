# Development TODO

## Task Completion Process

1. Move completed tasks to "Ready for Review" section
2. Mark with `[x]` when fully implemented AND tested
3. No task is complete until tests pass

### 1. ⏳ Implement Vault Client

**File:** `internal/sm/vault.go`

Implement `SecretManager` interface for HashiCorp Vault KV v2.

- [ ] Add `github.com/hashicorp/vault/api` to go.mod
- [ ] Create `VaultClient` struct
- [ ] Implement `NewVaultClient(addr, token string)`
- [ ] Implement `GetKey(path string) (*KeyPair, error)`
- [ ] Implement `StoreKey(path string, key *KeyPair) error`
- [ ] Implement `Close() error`
- [ ] Handle KV v2 `data` wrapper
- [ ] Handle auth errors gracefully

---

### 2. ⏳ Wire CLI to cmd Package

**File:** `internal/cli/cli.go`

- [ ] Import `cmd` package
- [ ] `runGenerate()` → call `cmd.Generate()`
- [ ] `runLoad()` → call `cmd.Load()`
- [ ] `NewSecretManager()` → create Vault client

---

### 3. ⏳ Implement Entry Point

**File:** `main.go`

- [ ] Call `cli.Run()`
- [ ] Handle exit codes

---

### 4. ⏳ Add Mock for Testing

**File:** `internal/sm/mock_vault.go`

- [ ] Create `MockSecretManager` struct
- [ ] Implement `GetKey(path string) (*KeyPair, error)` with in-memory map
- [ ] Implement `StoreKey(path string, key *KeyPair) error`
- [ ] Implement `Close() error` (no-op)
- [ ] Add `NewMockSecretManager()` constructor

---

### 5. ⏳ Add Tests

**Files:** `internal/sm/vault_test.go`, `internal/sm/mock_vault_test.go`

- [ ] Test `StoreKey()` and `GetKey()` with mock
- [ ] Test error handling (key not found, connection errors)
- [ ] Test KV v2 data wrapper serialization

---

## Testing

```bash
go test ./...
```
