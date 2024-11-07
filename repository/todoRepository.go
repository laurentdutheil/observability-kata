package repository

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"todo_odd/domain"
)

type TodoRepository struct {
	todos []domain.Todo
}

func (r *TodoRepository) AddTodo(ctx context.Context, title string, description string) (domain.Todo, error) {
	instrumentation := startInstrumentation(ctx, "todo creation repo")
	defer instrumentation.stopInstrumentation()

	todo := domain.Todo{
		Id:          len(r.todos) + 1,
		Title:       title,
		Description: description,
	}
	r.todos = append(r.todos, todo)

	instrumentation.todoCreated(todo.Id)

	return todo, nil
}

func (r *TodoRepository) Get(id int) (domain.Todo, error) {
	if id > len(r.todos) {
		return domain.Todo{}, fmt.Errorf("repository: todo #%d does not exist", id)
	}
	return r.todos[id-1], nil
}

func (r *TodoRepository) All() ([]domain.Todo, error) {
	return r.todos, nil
}

type TodoRepositoryInstrumentation struct {
	ctx    context.Context
	tracer trace.Tracer
	span   trace.Span
}

func startInstrumentation(ctx context.Context, name string) *TodoRepositoryInstrumentation {
	i := &TodoRepositoryInstrumentation{ctx: ctx, tracer: otel.Tracer("")}
	i.ctx, i.span = i.tracer.Start(i.ctx, name)
	return i
}

func (i *TodoRepositoryInstrumentation) stopInstrumentation() {
	i.span.End()
}

func (i *TodoRepositoryInstrumentation) todoCreated(id int) {
	i.span.SetAttributes(attribute.Int("id", id))
}
