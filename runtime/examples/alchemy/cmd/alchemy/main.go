package main

import (
	"fmt"
	"io"
	"os"

	"github.com/meowpow-png/mipe/runtime/examples/alchemy/internal/output"
	"github.com/meowpow-png/mipe/runtime/examples/alchemy/internal/storage"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return usageError()
	}
	book, err := storage.Load("potion.config")
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		if len(args) != 1 {
			return usageError()
		}
		output.List(stdout, book.All())
	case "show", "brew":
		if len(args) != 2 {
			return usageError()
		}
		item, err := book.Find(args[1])
		if err != nil {
			return err
		}
		if args[0] == "show" {
			output.Show(stdout, item)
		} else {
			output.Brew(stdout, item)
		}
	default:
		return usageError()
	}
	return nil
}

func usageError() error {
	return fmt.Errorf("usage: alchemy list | alchemy show <name> | alchemy brew <name>")
}
