package domain

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoService struct {
	Repository TodoRepository
}

func (s TodoService) AddTodo(ctx context.Context, title string, description string) Todo {
	todo := s.Repository.AddTodo(ctx, title, description)

	span := trace.SpanFromContext(ctx)
	span.SetName("todo creation")
	span.SetAttributes(attribute.Int("id", todo.Id))
	return todo
}

func (s TodoService) GetTodo(todoId int) Todo {
	return s.Repository.Get(todoId)
}

func (s TodoService) GetAll() []Todo {
	return s.Repository.All()
}

func (s TodoService) AddAllTodos(ctx context.Context, requestTodos []Todo) []Todo {
	span := trace.SpanFromContext(ctx)
	span.SetName("todo creation all")

	var todos []Todo
	for _, requestTodo := range requestTodos {
		todo := s.Repository.AddTodo(ctx, requestTodo.Title, requestTodo.Description)
		todos = append(todos, todo)
	}
	return todos
}
