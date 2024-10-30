package repository

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"todo_odd/domain"
)

type TodoRepository struct {
	todos []domain.Todo
}

var tracer = otel.Tracer("")

func (r *TodoRepository) AddTodo(ctx context.Context, title string, description string) domain.Todo {
	_, span := tracer.Start(ctx, "repository creation")
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

func (r *TodoRepository) Get(id int) (domain.Todo, error) {
	if id > len(r.todos) {
		return domain.Todo{}, fmt.Errorf("repository: todo #%d does not exist", id)
	}
	return r.todos[id-1], nil
}

func (r *TodoRepository) All() []domain.Todo {
	return r.todos
}
