package repository

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Todo struct {
	Id          int
	Title       string
	Description string
}

type TodoRepository struct {
	todos []Todo
}

func (r *TodoRepository) AddTodo(ctx context.Context, title string, description string) Todo {
	_, span := otel.Tracer("").Start(ctx, "repository creation")
	defer span.End()

	todo := Todo{
		Id:          len(r.todos) + 1,
		Title:       title,
		Description: description,
	}

	span.SetName("todo creation repo")
	span.SetAttributes(attribute.Int("id", todo.Id))

	r.todos = append(r.todos, todo)
	return todo
}

func (r *TodoRepository) Get(id int) Todo {
	return r.todos[id-1]
}

func (r *TodoRepository) All() []Todo {
	return r.todos
}
