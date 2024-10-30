package domain

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type TodoService struct {
	Repository TodoRepository
}

var tracer = otel.Tracer("")

func (s TodoService) AddTodo(ctx context.Context, title string, description string) Todo {
	ctx, span := tracer.Start(ctx, "todo creation")
	defer span.End()

	todo := s.Repository.AddTodo(ctx, title, description)

	span.SetAttributes(attribute.Int("id", todo.Id))

	return todo
}

func (s TodoService) GetTodo(todoId int) (Todo, error) {
	todo, err := s.Repository.Get(todoId)
	if err != nil {
		return Todo{}, fmt.Errorf("todo #%d does not exist", todoId)
	}
	return todo, err
}

func (s TodoService) GetAll() []Todo {
	return s.Repository.All()
}

func (s TodoService) AddAllTodos(ctx context.Context, requestTodos []Todo) []Todo {
	ctx, span := tracer.Start(ctx, "todo creation all")
	defer span.End()

	var todos []Todo
	for _, requestTodo := range requestTodos {
		todo := s.Repository.AddTodo(ctx, requestTodo.Title, requestTodo.Description)
		todos = append(todos, todo)
	}
	return todos
}
