# Contributing

Thanks for your interest in contributing!

## Quick Start

1. Fork the repository
2. Create a branch: `git checkout -b my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit: `git commit -m "feat: add my feature"`
6. Push: `git push origin my-feature`
7. Open a pull request

## Development

```bash
# Install dependencies
go mod download

# Format code
gofmt -w .
```

## Commit Messages

We use conventional commits:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `test:` - Tests
- `ci:` - CI/CD
- `chore:` - Maintenance

## Adding a New Secret Manager Provider

To add support for a new secret manager (e.g., AWS Secrets Manager, Azure Key Vault):

### 1. Add Provider Constant

In `internal/sm/sm.go`, add your provider constant:
```go
const (
    ProviderVault = "vault"
    ProviderAWS   = "aws"  // Add this
)
```

### 2. Create Provider Client

Create a new file in `internal/sm/` (e.g., `aws.go`) with your client implementation:
```go
type AWSClient struct {
    // your fields
}

func NewAWSClient(cfg *config.Config) (*AWSClient, error) {
    // initialization
}

func (a *AWSClient) GetKV(path string) (*KeyValue, error) {
    // implementation
}

func (a *AWSClient) StoreKV(path string, kv *KeyValue) error {
    // implementation
}

func (a *AWSClient) CheckExists(path string) (bool, error) {
    // implementation
}
```

Reference: `internal/sm/vault.go` for complete example.

### 3. Update sm.go Functions

In `internal/sm/sm.go`, add cases for your provider in `Get()`, `Store()`, and `CheckExists()`:

```go
case ProviderAWS:
    client, err := NewAWSClient(cfg)
    if err != nil {
        return false, err  // or return nil, err for Get()
    }
    return client.CheckExists(path)  // or GetKV(path) or StoreKV(path, kv)
```

### 4. Update Config

In `internal/config/config.go`:
1. Add provider paths field to `Config` struct:
```go
type Config struct {
    DefaultProvider string   `json:"default_provider"`
    VaultPaths      []string `json:"vault_paths,omitempty"`
    AWSPaths        []string `json:"aws_paths,omitempty"`  // Add this
}
```

2. Add validation in `Read()` function:
```go
switch cfg.DefaultProvider {
case sm.ProviderVault:
case sm.ProviderAWS:  // Add this
default:
    return nil, ErrInvalidProvider
}
```

3. Add getter method:
```go
func (c *Config) GetAWSPaths() []string {
    if c.AWSPaths == nil {
        return []string{}
    }
    return c.AWSPaths
}
```

### 5. Add Tests

Create both unit and integration tests:

**Unit test** in `internal/sm/aws_test.go`:
```go
package sm

import "testing"

func TestNewAWSClient(t *testing.T) {
    // test implementation
}
```

**Integration test** in `internal/sm/aws_integration_test.go`:
```go
//go:build integration

package sm

import (
    "testing"
)

func TestAWSClientGetKV(t *testing.T) {
    // test with real AWS service
}
```

### 6. Update CI Workflows

Add your provider to the integration test matrix in `.github/workflows/tests-integration.yml`:

```yaml
strategy:
  matrix:
    include:
      - provider: Vault
        image: hashicorp/vault:latest
        addr_env: VAULT_ADDR
        token_env: VAULT_TOKEN
        service_env_prefix: VAULT
      - provider: AWS  # Add this
        image: localstack/localstack:latest
        addr_env: AWS_ENDPOINT_URL
        token_env: AWS_ACCESS_KEY_ID
        service_env_prefix: AWS
```

### 7. Update Documentation

Add your provider to:
- `README.md` - Supported providers list, environment variables, configuration examples
- `internal/sm/sm.go` - Provider constants comment

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
