package store

import (
	"encoding/json"
	"os"

	"go-cli-todo-list/internal/todo"
)

const defaultFile = "todos.json"

type Store struct {
	filepath string
}

func New(filepath string) *Store {
	if filepath == "" {
		filepath = defaultFile
	}
	return &Store{filepath: filepath}
}

func (s *Store) Load() (todo.List, error) {
	var list todo.List
	data, err := os.ReadFile(s.filepath)
	if os.IsNotExist(err) {
		return list, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Save(list todo.List) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filepath, data, 0644)
}
