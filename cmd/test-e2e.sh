#!/bin/bash
# End-to-End Integration Test for sm-ssh-add
# This script tests the complete CLI workflow by building and running the actual binary
#
# Environment Variables:
#   PROVIDER - The secret manager provider to test (vault|openbao)
#   VAULT_ADDR / VAULT_TOKEN - For Vault provider
#   BAO_ADDR / BAO_TOKEN - For OpenBao provider
#   SSH_AUTH_SOCK - Required for ssh-agent tests
#
# Usage:
#   ./cmd/test-e2e.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test configuration
PROVIDER="${PROVIDER:-vault}"
TEST_DIR=$(mktemp -d)
TEST_KEY_PATH_1="secret/data/ssh/e2e-test-key-1"
TEST_KEY_PATH_2="secret/data/ssh/e2e-test-key-2"
TEST_KEY_PATH_PASSPHRASE="secret/data/ssh/e2e-test-passphrase"
TEST_PASSPHRASE="test-passphrase-12345"
CONFIG_FILE="$TEST_DIR/sm-ssh-add.json"

# Cleanup on exit
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"
    # Note: We don't delete keys from vault as the environment is ephemeral
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Print header
print_header() {
    echo ""
    echo "========================================"
    echo "  E2E Test: $1"
    echo "  Provider: $PROVIDER"
    echo "========================================"
    echo ""
}

# Print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"

    # Check provider environment variables
    if [ "$PROVIDER" = "vault" ]; then
        [ -n "$VAULT_ADDR" ] || print_error "VAULT_ADDR not set"
        [ -n "$VAULT_TOKEN" ] || print_error "VAULT_TOKEN not set"
        print_success "Vault environment variables set"
    elif [ "$PROVIDER" = "openbao" ]; then
        [ -n "$BAO_ADDR" ] || print_error "BAO_ADDR not set"
        [ -n "$BAO_TOKEN" ] || print_error "BAO_TOKEN not set"
        print_success "OpenBao environment variables set"
    else
        print_error "Invalid provider: $PROVIDER (must be 'vault' or 'openbao')"
    fi

    # Check ssh-agent
    if [ -z "$SSH_AUTH_SOCK" ]; then
        print_error "SSH_AUTH_SOCK not set - ssh-agent not running"
    fi
    print_success "SSH_AUTH_SOCK is set"

    # Check go compiler
    command -v go >/dev/null 2>&1 || print_error "Go compiler not found"
    print_success "Go compiler available"
}

# Build the binary
build_binary() {
    print_header "Building Binary"

    cd "$(dirname "$0")/.."
    go build -o "$TEST_DIR/sm-ssh-add" . || print_error "Build failed"
    print_success "Binary built: $TEST_DIR/sm-ssh-add"
    BINARY="$TEST_DIR/sm-ssh-add"
}

# Create test config
create_config() {
    cat > "$CONFIG_FILE" <<EOF
{
  "default_provider": "$PROVIDER",
  "vault_paths": []
}
EOF
    # Set up config in test home directory
    export HOME="$TEST_DIR"
    mkdir -p "$TEST_DIR/.config"
    cp "$CONFIG_FILE" "$TEST_DIR/.config/sm-ssh-add.json"
    print_success "Config created with provider: $PROVIDER"
}

# Test 1: Generate and load key without passphrase
test_generate_load_without_passphrase() {
    print_header "Test 1: Generate & Load (No Passphrase)"

    # Generate key
    echo "Generating key at: $TEST_KEY_PATH_1"
    "$BINARY" generate "$TEST_KEY_PATH_1" "test1@example.com" > "$TEST_DIR/generate-output.txt" 2>&1
    print_success "Key generated"

    # Verify public key was output
    grep -q "ssh-ed25519" "$TEST_DIR/generate-output.txt" || print_error "Public key not found in output"
    print_success "Public key format verified"

    # Load key
    echo "Loading key from: $TEST_KEY_PATH_1"
    "$BINARY" load "$TEST_KEY_PATH_1" > "$TEST_DIR/load-output.txt" 2>&1
    print_success "Key loaded into ssh-agent"

    # Verify key is in agent
    ssh-add -l | grep -q "test1@example.com" || print_error "Key not found in ssh-agent"
    print_success "Key verified in ssh-agent"
}

# Test 2: Generate and load key with passphrase
test_generate_load_with_passphrase() {
    print_header "Test 2: Generate & Load (With Passphrase)"

    # Generate key with passphrase
    echo "Generating key with passphrase at: $TEST_KEY_PATH_PASSPHRASE"
    # Provide passphrase twice for confirmation during generate
    echo -e "$TEST_PASSPHRASE\n$TEST_PASSPHRASE" | "$BINARY" generate --require-passphrase "$TEST_KEY_PATH_PASSPHRASE" "test-pass@example.com" > "$TEST_DIR/generate-pass-output.txt" 2>&1
    print_success "Key with passphrase generated"

    # Verify public key was output
    grep -q "ssh-ed25519" "$TEST_DIR/generate-pass-output.txt" || print_error "Public key not found in output"
    print_success "Public key format verified"

    # Load key with passphrase
    echo "Loading passphrase-protected key"
    echo "$TEST_PASSPHRASE" | "$BINARY" load "$TEST_KEY_PATH_PASSPHRASE" > "$TEST_DIR/load-pass-output.txt" 2>&1
    print_success "Passphrase-protected key loaded into ssh-agent"

    # Verify key is in agent
    ssh-add -l | grep -q "test-pass@example.com" || print_error "Passphrase-protected key not found in ssh-agent"
    print_success "Passphrase-protected key verified in ssh-agent"
}

# Test 3: Load multiple keys from config
test_load_from_config() {
    print_header "Test 3: Load Multiple Keys from Config"

    # Generate second key
    echo "Generating second key at: $TEST_KEY_PATH_2"
    "$BINARY" generate "$TEST_KEY_PATH_2" "test2@example.com" > "$TEST_DIR/generate-output-2.txt" 2>&1
    print_success "Second key generated"

    # Update config to include both paths
    cat > "$TEST_DIR/.config/sm-ssh-add.json" <<EOF
{
  "default_provider": "$PROVIDER",
  "vault_paths": ["$TEST_KEY_PATH_1", "$TEST_KEY_PATH_2"]
}
EOF
    print_success "Config updated with multiple paths"

    # Remove all keys from agent first
    ssh-add -D > /dev/null 2>&1 || true

    # Load from config
    echo "Loading keys from config..."
    "$BINARY" load --from-config > "$TEST_DIR/load-config-output.txt" 2>&1
    print_success "Keys loaded from config"

    # Verify both keys are in agent
    ssh-add -l | grep -q "test1@example.com" || print_error "Key 1 not found in ssh-agent"
    ssh-add -l | grep -q "test2@example.com" || print_error "Key 2 not found in ssh-agent"
    print_success "Both keys verified in ssh-agent"
}

# Test 4: Duplicate key detection
test_duplicate_detection() {
    print_header "Test 4: Duplicate Key Detection"

    # Try to load the same key again
    echo "Attempting to load duplicate key..."
    output=$("$BINARY" load "$TEST_KEY_PATH_1" 2>&1)

    # Should indicate key already loaded
    if echo "$output" | grep -q "already loaded"; then
        print_success "Duplicate key detected and skipped"
    else
        print_error "Duplicate key was not detected"
    fi
}

# Run all tests
main() {
    echo ""
    echo "========================================"
    echo "  sm-ssh-add End-to-End Integration Test"
    echo "========================================"
    echo ""

    check_prerequisites
    build_binary
    create_config

    # Run tests
    test_generate_load_without_passphrase
    test_generate_load_with_passphrase
    test_load_from_config
    test_duplicate_detection

    # Final summary
    print_header "All Tests Passed!"
    echo -e "${GREEN}✓ Generate & load (no passphrase)${NC}"
    echo -e "${GREEN}✓ Generate & load (with passphrase)${NC}"
    echo -e "${GREEN}✓ Load multiple keys from config${NC}"
    echo -e "${GREEN}✓ Duplicate key detection${NC}"
    echo ""
    echo -e "${GREEN}========================================"
    echo "  SUCCESS: All E2E tests passed!"
    echo "========================================${NC}"
}

main
