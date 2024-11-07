package domain

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoService struct {
	Repository TodoRepository
}

func (s TodoService) AddTodo(ctx context.Context, title string, description string) Todo {
	instrumentation := startInstrumentation(ctx, "todo creation")
	defer instrumentation.stopInstrumentation()

	todo, _ := s.Repository.AddTodo(instrumentation.ctx, title, description)

	instrumentation.todoCreated(todo.Id)

	return todo
}

func (s TodoService) GetTodo(todoId int) (Todo, error) {
	todo, err := s.Repository.Get(todoId)
	if err != nil {
		return Todo{}, fmt.Errorf("todo #%d does not exist", todoId)
	}
	return todo, err
}

func (s TodoService) GetAll() ([]Todo, error) {
	return s.Repository.All()
}

func (s TodoService) AddAllTodos(ctx context.Context, requestTodos []Todo) []Todo {
	instrumentation := startInstrumentation(ctx, "todo creation all")
	defer instrumentation.stopInstrumentation()

	var todos []Todo
	for _, requestTodo := range requestTodos {
		todo, _ := s.Repository.AddTodo(instrumentation.ctx, requestTodo.Title, requestTodo.Description)
		todos = append(todos, todo)
	}
	return todos
}

type ServiceInstrumentation struct {
	ctx    context.Context
	tracer trace.Tracer
	span   trace.Span
}

func startInstrumentation(ctx context.Context, name string) *ServiceInstrumentation {
	i := &ServiceInstrumentation{ctx: ctx, tracer: otel.Tracer("")}
	i.ctx, i.span = i.tracer.Start(i.ctx, name)
	return i
}

func (i *ServiceInstrumentation) stopInstrumentation() {
	i.span.End()
}

func (i *ServiceInstrumentation) todoCreated(id int) {
	i.span.SetAttributes(attribute.Int("id", id))
}
