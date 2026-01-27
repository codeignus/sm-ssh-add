package sm

import (
	"errors"
	"fmt"
)

// Errors returned by the secret manager package
var (
	ErrPathNotFound     = errors.New("path not found in secret manager")
	ErrInvalidKeyFormat = errors.New("invalid key format in secret manager")
	ErrKeyExistsInAgent = errors.New("key already exists in ssh-agent")
	ErrVaultConnection  = errors.New("failed to connect to vault")
	ErrSSHAgentNotFound = errors.New("ssh-agent not found")
)

// wrapError wraps an error with additional context
func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
