package main

import (
	"fmt"
	"os"

	"github.com/codeignus/sm-ssh-add/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: sm-ssh-add <generate>\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		if err := cmd.Generate(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "usage: sm-ssh-add <generate>\n")
		os.Exit(1)
	}
}
