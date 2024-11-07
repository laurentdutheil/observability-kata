package repository

import (
	"context"
	"fmt"
	"todo_odd/domain"
)

type InMemoryRepository struct {
	todos []domain.Todo
}

func (r *InMemoryRepository) AddTodo(ctx context.Context, title string, description string) (domain.Todo, error) {
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

func (r *InMemoryRepository) Get(id int) (domain.Todo, error) {
	if id == 0 || id > len(r.todos) {
		return domain.Todo{}, fmt.Errorf("repository: todo #%d does not exist", id)
	}
	return r.todos[id-1], nil
}

func (r *InMemoryRepository) All() ([]domain.Todo, error) {
	return r.todos, nil
}
