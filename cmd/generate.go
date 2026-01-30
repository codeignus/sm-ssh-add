package cmd

import (
	"fmt"
	"os"

	"github.com/codeignus/sm-ssh-add/internal/config"
	"github.com/codeignus/sm-ssh-add/internal/sm"
	"github.com/codeignus/sm-ssh-add/internal/ssh"
)

// Generate creates a new SSH key pair and displays the public key
func Generate(cfg *config.Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: sm-ssh-add generate [--require-passphrase] <path> [comment]")
	}

	if len(args) > 3 {
		return fmt.Errorf("too many arguments\nusage: sm-ssh-add generate [--require-passphrase] <path> [comment]")
	}

	// Parse arguments
	path := ""
	comment := ""
	requirePassphrase := false

	for _, arg := range args {
		switch arg {
		case "--require-passphrase":
			requirePassphrase = true
		default:
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

	err = sm.Store(cfg.DefaultProvider, path, kv)
	if err != nil {
		return fmt.Errorf("failed to store key in vault: %w", err)
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
