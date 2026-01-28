package ssh

import (
	"fmt"
)

// wrapError wraps an error with context
func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
