package config

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrConfigFileNotFound = errors.New("config file not found")
	ErrEmptyProvider      = errors.New("default_provider cannot be empty")
	ErrInvalidProvider    = errors.New("invalid provider")
)

// wrapError wraps an error with context
func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
