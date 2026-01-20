package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Generate creates a new SSH key pair and displays the public key
func Generate() error {
	requirePassphrase := flag.Bool("require-passphrase", false, "prompt for passphrase to protect the private key")
	flag.CommandLine.Parse(os.Args[2:])

	args := flag.CommandLine.Args()

	if len(args) == 0 {
		return fmt.Errorf("usage: sm-ssh-add generate [--require-passphrase] <path> [comment]")
	}

	if len(args) > 2 {
		return fmt.Errorf("too many arguments\nusage: sm-ssh-add generate [--require-passphrase] <path> [comment]")
	}

	path := args[0]
	comment := ""
	if len(args) == 2 {
		comment = args[1]
	}

	// TODO: GenerateKeyPair

	// TODO: store key in vault

	fmt.Fprintf(os.Stdout, "Generated SSH key for path: %s\n", path)

	return nil
}
