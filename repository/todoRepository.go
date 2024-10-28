package repository

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"todo_odd/domain"
)

type TodoRepository struct {
	todos []domain.Todo
}

func (r *TodoRepository) AddTodo(ctx context.Context, title string, description string) domain.Todo {
	_, span := otel.Tracer("").Start(ctx, "repository creation")
	defer span.End()

	todo := domain.Todo{
		Id:          len(r.todos) + 1,
		Title:       title,
		Description: description,
	}

	span.SetName("todo creation repo")
	span.SetAttributes(attribute.Int("id", todo.Id))

	r.todos = append(r.todos, todo)
	return todo
}

func (r *TodoRepository) Get(id int) domain.Todo {
	return r.todos[id-1]
}

func (r *TodoRepository) All() []domain.Todo {
	return r.todos
}
