package cmd

import (
	"fmt"
	"os"

	"github.com/codeignus/sm-ssh-add/internal/config"
	"github.com/codeignus/sm-ssh-add/internal/sm"
	"github.com/codeignus/sm-ssh-add/internal/ssh"
)

// Generate creates a new SSH key pair and displays the public key
func Generate(provider sm.Provider, cfg *config.Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: sm-ssh-add generate [--require-passphrase] [--save-path] <path> [comment]")
	}

	if len(args) > 5 {
		return fmt.Errorf("too many arguments\nusage: sm-ssh-add generate [--require-passphrase] [--save-path] [--regenerate] <path> [comment]")
	}

	// Parse arguments
	path := ""
	comment := ""
	requirePassphrase := false
	savePath := false
	regenerateKeypair := false

	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			switch arg {
			case "--require-passphrase":
				requirePassphrase = true
			case "--save-path":
				savePath = true
			case "--regenerate":
				regenerateKeypair = true
			default:
				return fmt.Errorf("unknown flag: %s", arg)
			}
		} else {
			if path == "" {
				path = arg
			} else if comment == "" {
				comment = arg
			}
		}
	}

	if path == "" {
		return fmt.Errorf("path is required\nusage: sm-ssh-add generate [--require-passphrase] <path> [comment]")
	}

	var passphrase []byte
	var err error

	// If not regenerating, check if key already exists
	if !regenerateKeypair {
		exists, err := provider.CheckExists(path)
		if err != nil {
			return fmt.Errorf("failed to check existing key: %w", err)
		}
		if exists {
			return fmt.Errorf("key already exists at %s (use --regenerate to overwrite)", path)
		}
	}

	if requirePassphrase {
		passphrase, err = readPassphrase()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}
	}

	// Generate the key pair
	keyPair, err := ssh.GenerateKeyPair(comment, passphrase)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Store in vault
	kv := &sm.KeyValue{
		PrivateKey:        keyPair.PrivateKey,
		PublicKey:         keyPair.PublicKey,
		RequirePassphrase: requirePassphrase,
		Comment:           comment,
	}

	err = provider.Store(path, kv)
	if err != nil {
		return fmt.Errorf("failed to store key in vault: %w", err)
	}

	// Save path to config if requested
	if savePath {
		if err := cfg.AddPath(path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Path %q saved to config\n", path)
		}
	}

	// Display the public key
	fmt.Fprintf(os.Stdout, "%s", keyPair.PublicKey)
	fmt.Fprintf(os.Stdout, "Key stored at: %s\n", path)

	return nil
}

// readPassphrase reads a passphrase from stdin twice to confirm
func readPassphrase() ([]byte, error) {
	fmt.Fprint(os.Stderr, "Enter passphrase (empty for no passphrase): ")
	var passphrase1 string
	_, err := fmt.Scanln(&passphrase1)
	if err != nil {
		return nil, err
	}

	fmt.Fprint(os.Stderr, "Enter same passphrase again: ")
	var passphrase2 string
	_, err = fmt.Scanln(&passphrase2)
	if err != nil {
		return nil, err
	}

	if passphrase1 != passphrase2 {
		return nil, fmt.Errorf("passphrases do not match")
	}

	return []byte(passphrase1), nil
}
