# sm-ssh-add

A CLI tool that generates SSH keys, stores them in Secret Managers like HashiCorp Vault, OpenBao etc and loads them into ssh-agent.

**Goal:** Eliminate local SSH key storage while maintaining seamless ssh-agent integration.

## Features

- Ed25519 SSH key generation with optional passphrase protection
- Vault KV v2 storage for secure key management
- ssh-agent integration with duplicate detection
- Multi-key loading from configured paths
- JSON configuration for flexible setup

## Installation

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/codeignus/sm-ssh-add/releases).

```bash
# Linux (x86_64)
curl -LO https://github.com/codeignus/sm-ssh-add/releases/latest/download/sm-ssh-add_0.2.0_Linux_x86_64.tar.gz
mkdir -p /tmp/sm-ssh-add && tar -xzf sm-ssh-add_0.2.0_Linux_x86_64.tar.gz -C /tmp/sm-ssh-add
sudo mv /tmp/sm-ssh-add/sm-ssh-add /usr/local/bin/
rm -rf /tmp/sm-ssh-add sm-ssh-add_0.2.0_Linux_x86_64.tar.gz

# macOS (Apple Silicon)
curl -LO https://github.com/codeignus/sm-ssh-add/releases/latest/download/sm-ssh-add_0.2.0_Darwin_arm64.tar.gz
mkdir -p /tmp/sm-ssh-add && tar -xzf sm-ssh-add_0.2.0_Darwin_arm64.tar.gz -C /tmp/sm-ssh-add
sudo mv /tmp/sm-ssh-add/sm-ssh-add /usr/local/bin/
rm -rf /tmp/sm-ssh-add sm-ssh-add_0.2.0_Darwin_arm64.tar.gz
```

### Build from Source

```bash
git clone https://github.com/codeignus/sm-ssh-add.git
cd sm-ssh-add
go build -o sm-ssh-add .
sudo mv sm-ssh-add /usr/local/bin/
```

## Usage

### Prerequisites

Before using `sm-ssh-add`, create a configuration file:

**Location:** `~/.config/sm-ssh-add.json`

```json
{
  "default_provider": "vault",
  "vault_paths": []
}
```

## Commands

### generate

Create a new SSH key pair and store it in Vault.

```bash
sm-ssh-add generate [--require-passphrase] <path> [comment]
```

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `path` | ✅ | Path to your configured secret manager (e.g., `secret/ssh/github`) |
| `comment` | ❌ | Key comment (email or identifier) |

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--require-passphrase` | `false` | Prompt for passphrase |

**Examples:**

```bash
sm-ssh-add generate secret/ssh/github "user@example.com"
sm-ssh-add generate --require-passphrase secret/ssh/gitlab "user@example.com"
```

**Output:**

```
Generated SSH key for path: secret/ssh/github
Public key: ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExampleKey user@example.com
```

### load

Load SSH keys from Vault into ssh-agent.

```bash
# Load all keys from configuration file
sm-ssh-add load --from-config

# Load a specific key from a path to your configured secret manager
sm-ssh-add load secret/ssh/github
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--from-config` | Load all keys from the configured `<provider>_paths` in the config file |

**Arguments:**

When not using `--from-config`, you must provide a path to your configured secret manager as an argument.

| Argument | Required | Description |
|----------|----------|-------------|
| `path` | ❌ | Path to your configured secret manager (required only when not using `--from-config`) |

**Examples:**

```bash
# Load all configured keys
sm-ssh-add load --from-config

# Load a single key
sm-ssh-add load secret/ssh/github
```

**Output:**

```
Loaded key: github (ssh-ed25519 AAAAC3NzaC1lZDI1NTE5...)
Loaded key: gitlab (ssh-ed25519 AAAAC3NzaC1lZDI1NTE5...)
Skipped: staging (already in agent)
✓ Loaded 2 keys
```

## Configuration

The configuration file must be created at `~/.config/sm-ssh-add.json` before running any commands (see [Usage](#usage) above).

### Configuration Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `default_provider` | string | ✅ | Secret manager to use ("vault" only) |
| `<provider>_paths` | string[] | ✅ | List of paths to load keys from (e.g., `vault_paths`, `aws_paths`) |
| `vault_approle_role_id` | string | ❌ | AppRole Role ID for Vault auth (uses AppRole instead of VAULT_TOKEN; Secret ID via VAULT_APPROLE_SECRET_ID or prompt) |

### Environment Variables

| Variable | Used For |
|----------|----------|
| `HOME` | Locating config file |
| `SSH_AUTH_SOCK` | Default SSH agent socket (fallback) |
| `VAULT_ADDR` / `BAO_ADDR` | Vault or OpenBao server address |
| `VAULT_TOKEN` / `BAO_TOKEN` | Authentication token (used when `vault_approle_role_id` is not set) |
| `VAULT_APPROLE_SECRET_ID` | AppRole secret ID (optional when `vault_approle_role_id` is set; prompted if not provided) |
