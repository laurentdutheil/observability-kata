package domain

import (
	"context"
)

type TodoRepository interface {
	AddTodo(ctx context.Context, title string, description string) Todo
	Get(id int) Todo
	All() []Todo
}
