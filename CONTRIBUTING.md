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
- `refactor:` - Code refactoring
- `docs:` - Documentation
- `test:` - Tests
- `ci:` - CI/CD
- `chore:` - Maintenance

## Adding a New Secret Manager Provider

To add support for a new secret manager (e.g., AWS Secrets Manager, Azure Key Vault):

### 1. Add Provider Constant

In `internal/config/config.go`, add your provider constant:
```go
const (
    ProviderVault = "vault"
    ProviderAWS   = "aws"  // Add this
)
```

### 2. Create Provider Client

Create a new file in `internal/sm/` (e.g., `aws.go`) with your client implementation that directly implements the `Provider` interface:

```go
package sm

import "github.com/codeignus/sm-ssh-add/internal/config"

// AWSClient implements secret manager operations for AWS Secrets Manager
type AWSClient struct {
    // your fields
}

// NewAWSClient creates a new AWS client
func NewAWSClient(cfg *config.Config) (*AWSClient, error) {
    // initialization
}

// Get retrieves key-value data from AWS
// Implements Provider interface
func (a *AWSClient) Get(path string) (*KeyValue, error) {
    // implementation - convert AWS response to KeyValue format
}

// Store stores key-value data to AWS
// Implements Provider interface
func (a *AWSClient) Store(path string, kv *KeyValue) error {
    // implementation - convert KeyValue to AWS format
}

// CheckExists checks if a key exists at the given path
// Implements Provider interface
func (a *AWSClient) CheckExists(path string) (bool, error) {
    // implementation
}
```

Reference: `internal/sm/vault.go` for a complete example.

### 3. Update InitProvider

In `internal/sm/sm.go`, update the `InitProvider()` function to return your new client:

```go
func InitProvider(cfg *config.Config) (Provider, error) {
    switch cfg.DefaultProvider {
    case config.ProviderVault:
        return NewVaultClient(cfg)
    case config.ProviderAWS:  // Add this
        return NewAWSClient(cfg)
    // Future: add cases for Azure, GCP, etc.
    default:
        return nil, fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
    }
}
```

### 4. Update Config

In `internal/config/config.go`:

1. Add provider paths field to `Config` struct:
```go
type Config struct {
    DefaultProvider string   `json:"default_provider"`
    VaultPaths      []string `json:"vault_paths,omitempty"`
    AWSPaths        []string `json:"aws_paths,omitempty"`  // Add this
    VaultApproleRoleID string `json:"vault_approle_role_id,omitempty"`
}
```

2. Add validation in `Read()` function:
```go
switch cfg.DefaultProvider {
case ProviderVault:
case ProviderAWS:  // Add this
default:
    return nil, ErrInvalidProvider
}
```

3. Add cases to `GetPaths()` and `AddPath()`:
```go
func (c *Config) GetPaths() []string {
    switch c.DefaultProvider {
    case ProviderVault:
        if c.VaultPaths == nil {
            return []string{}
        }
        return c.VaultPaths
    case ProviderAWS:  // Add this
        if c.AWSPaths == nil {
            return []string{}
        }
        return c.AWSPaths
    default:
        return []string{}
    }
}

func (c *Config) AddPath(path string) error {
    switch c.DefaultProvider {
    case ProviderVault:
        // ... existing logic
    case ProviderAWS:  // Add this
        if slices.Contains(c.AWSPaths, path) {
            return nil
        }
        c.AWSPaths = append(c.AWSPaths, path)
    default:
        return fmt.Errorf("unsupported provider: %s", c.DefaultProvider)
    }
    // ... write to disk logic
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
    "github.com/codeignus/sm-ssh-add/internal/config"
)

func TestAWSClientGet(t *testing.T) {
    cfg := &config.Config{DefaultProvider: config.ProviderAWS}
    client, err := NewAWSClient(cfg)
    if err != nil {
        t.Fatalf("Failed to create AWS client: %v", err)
    }
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
- This file (CONTRIBUTING.md) - Update the provider list in step 1

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
