package cmd

import (
	"fmt"
	"os"

	"github.com/codeignus/sm-ssh-add/internal/config"
	"github.com/codeignus/sm-ssh-add/internal/sm"
	"github.com/codeignus/sm-ssh-add/internal/ssh"
)

// parseLoadArgs parses command line arguments and returns the paths to load
func parseLoadArgs(args []string, cfg *config.Config) ([]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("usage: sm-ssh-add load [--from-config] <path>")
	}

	// Check for --from-config flag
	if args[0] == "--from-config" {
		if len(args) > 1 {
			return nil, fmt.Errorf("cannot use both --from-config and direct path")
		}
		paths := cfg.GetVaultPaths()
		if len(paths) == 0 {
			return nil, fmt.Errorf("no vault paths configured")
		}
		return paths, nil
	}

	// Check for unknown flags
	if args[0][0] == '-' {
		return nil, fmt.Errorf("unknown flag: %s", args[0])
	}

	// Direct path argument
	return []string{args[0]}, nil
}

// loadAndAddKey loads a key from the given path and adds it to the agent
func loadAndAddKey(path string, cfg *config.Config, agent *ssh.Agent) error {
	keyValue, err := sm.Get(cfg.DefaultProvider, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load key from %s: %v\n", path, err)
		return err
	}

	keyPair := &ssh.KeyPair{
		PrivateKey: keyValue.PrivateKey,
		PublicKey:  keyValue.PublicKey,
	}

	err = agent.AddKey(keyPair)
	if err != nil {
		if err == sm.ErrKeyExistsInAgent {
			fmt.Fprintf(os.Stdout, "Key from %s already loaded in agent\n", path)
			return nil
		}
		fmt.Fprintf(os.Stderr, "Failed to add key from %s: %v\n", path, err)
		return err
	}

	fmt.Fprintf(os.Stdout, "Loaded key from %s into ssh-agent\n", path)
	return nil
}

// Load retrieves SSH keys from the secret manager and adds them to ssh-agent
func Load(cfg *config.Config, args []string) error {
	paths, err := parseLoadArgs(args, cfg)
	if err != nil {
		return err
	}

	agent, err := ssh.NewAgent(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to ssh-agent: %w", err)
	}
	defer func() {
		if cerr := agent.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close ssh-agent: %v\n", cerr)
		}
	}()

	for _, path := range paths {
		if err := loadAndAddKey(path, cfg, agent); err != nil {
			return err
		}
	}

	return nil
}
