package domain

import (
	"context"
)

type TodoRepository interface {
	AddTodo(ctx context.Context, title string, description string) (Todo, error)
	Get(id int) (Todo, error)
	All() ([]Todo, error)
}
