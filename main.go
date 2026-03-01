package main

import (
	"fmt"
	"os"

	"go-cli-todo-list/internal/commands"
	"go-cli-todo-list/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		commands.New(store.New("")).Help()
		os.Exit(1)
	}

	s := store.New("")
	c := commands.New(s)

	cmd := os.Args[1]

	var err error
	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "error: 'add' requires a description\nUsage: todo add \"<description>\"")
			os.Exit(1)
		}
		err = c.Add(os.Args[2])

	case "list":
		err = c.List()

	case "done":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "error: 'done' requires an id\nUsage: todo done <id>")
			os.Exit(1)
		}
		err = c.Complete(os.Args[2])

	case "delete":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "error: 'delete' requires an id\nUsage: todo delete <id>")
			os.Exit(1)
		}
		err = c.Delete(os.Args[2])

	case "help":
		c.Help()

	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n%s", cmd, commands.Usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
