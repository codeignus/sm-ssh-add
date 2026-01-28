# Project Context

## File Structure

```
sm-ssh-add/
├── main.go                 # Entry point with CLI routing
├── cmd/                    # CLI commands
│   ├── generate.go         # Generate SSH key → Secret Manager
│   └── load.go             # Load keys from Secret Manager → ssh-agent
├── internal/
│   ├── sm/                 # Secret Manager package
│   │   ├── sm.go           # KeyValue struct, Get/Store functions, provider constants
│   │   ├── sm_test.go      # Tests for sm package
│   │   ├── errors.go       # Custom error definitions
│   │   ├── vault.go        # Vault/OpenBao KV v2 client implementation
│   │   └── vault_test.go   # Vault client tests
│   ├── config/             # Config file management
│   │   ├── config.go       # Config struct, Read() with validation
│   │   ├── config_test.go  # Config tests
│   │   └── errors.go       # Config-specific errors
│   └── ssh/                # Key generation + ssh-agent operations
│       ├── keygen.go       # Ed25519 key generation with comment/passphrase
│       ├── keygen_test.go  # Key generation tests
│       ├── agent.go        # ssh-agent operations (for load command)
│       └── errors.go       # Error wrapping
├── go.mod                  # Go module definition
└── go.sum                  # Go module checksums
```

### Planned Structure (Not Yet Implemented)

```
└── .github/workflows/      # CI/CD
    └── test.yml            # Integration tests with Vault/OpenBao containers
```

## Tech Stack

- Go 1.21+
- `golang.org/x/crypto/ssh` - SSH key generation
- `github.com/hashicorp/vault/api` - Secret Manager client (Vault/OpenBao compatible)
- Config: `~/.config/sm-ssh-add.json`

## Design Choices

- **Algorithm:** Ed25519 only (no RSA/ECDSA)
- **Storage:** Secret Manager (currently Vault KV v2, extensible for AWS Secrets Manager, etc.)
- **Duplicates:** SHA256 fingerprint check before adding to agent
- **Errors:** Wrapped with context, printed to stderr

## Secret Manager Data

### KeyValue Fields

- `private_key` - SSH private key in OpenSSH format
- `public_key` - SSH public key
- `require_passphrase` - Boolean indicating if key requires passphrase
