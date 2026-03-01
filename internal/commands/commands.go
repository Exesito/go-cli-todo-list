package commands

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"go-cli-todo-list/internal/store"
	"go-cli-todo-list/internal/todo"
)

type Commander struct {
	store *store.Store
}

func New(s *store.Store) *Commander {
	return &Commander{store: s}
}

func (c *Commander) Add(description string) error {
	list, err := c.store.Load()
	if err != nil {
		return fmt.Errorf("loading todos: %w", err)
	}
	list.Add(description)
	if err := c.store.Save(list); err != nil {
		return fmt.Errorf("saving todos: %w", err)
	}
	fmt.Printf("Todo added (#%d): %s\n", list[len(list)-1].ID, description)
	return nil
}

func (c *Commander) List() error {
	list, err := c.store.Load()
	if err != nil {
		return fmt.Errorf("loading todos: %w", err)
	}
	if len(list) == 0 {
		fmt.Println("No todos found. Use 'add <description>' to create one.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tCREATED AT\tDESCRIPTION")
	fmt.Fprintln(w, "--\t------\t----------\t-----------")
	for _, t := range list {
		status := "[ ] pending"
		if t.Status == todo.StatusCompleted {
			status = "[x] done   "
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			t.ID, status, t.CreatedAt.Format("2006-01-02 15:04"), t.Description)
	}
	w.Flush()
	return nil
}

func (c *Commander) Complete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %q: must be an integer", idStr)
	}
	list, err := c.store.Load()
	if err != nil {
		return fmt.Errorf("loading todos: %w", err)
	}
	if !list.Complete(id) {
		return fmt.Errorf("todo with id %d not found", id)
	}
	if err := c.store.Save(list); err != nil {
		return fmt.Errorf("saving todos: %w", err)
	}
	fmt.Printf("Todo #%d marked as complete.\n", id)
	return nil
}

func (c *Commander) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %q: must be an integer", idStr)
	}
	list, err := c.store.Load()
	if err != nil {
		return fmt.Errorf("loading todos: %w", err)
	}
	if !list.Delete(id) {
		return fmt.Errorf("todo with id %d not found", id)
	}
	if err := c.store.Save(list); err != nil {
		return fmt.Errorf("saving todos: %w", err)
	}
	fmt.Printf("Todo #%d deleted.\n", id)
	return nil
}

const Usage = `go-cli-todo-list — manage your tasks from the terminal

Usage:
  todo <command> [arguments]

Commands:
  add <description>   Add a new todo
  list                List all todos
  done <id>           Mark a todo as complete
  delete <id>         Delete a todo
  help                Show this help message

Examples:
  todo add "Buy groceries"
  todo list
  todo done 1
  todo delete 2
`

func (c *Commander) Help() {
	fmt.Print(Usage)
}
