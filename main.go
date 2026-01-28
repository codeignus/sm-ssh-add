package main

import (
	"fmt"
	"os"

	"github.com/codeignus/sm-ssh-add/cmd"
	"github.com/codeignus/sm-ssh-add/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: sm-ssh-add <generate|load> [args]\n")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "generate":
		if err := cmd.Generate(cfg, args); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "load":
		if err := cmd.Load(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "usage: sm-ssh-add <generate|load> [args]\n")
		os.Exit(1)
	}
}
