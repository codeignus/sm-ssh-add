# Project Context

## File Structure

```
sm-ssh-add/
├── main.go                 # Entry point, routes to commands
├── cmd/                    # CLI commands
│   ├── generate.go         # Generate SSH key → Secret Manager
│   └── load.go             # Load keys from Secret Manager → ssh-agent
├── internal/
│   ├── config/             # Config file management
│   ├── sm/                 # Secret Manager interface + implementations
│   └── ssh/                # Key generation + ssh-agent operations
└── go.mod
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
