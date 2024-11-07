package domain

import (
	"fmt"
)

type TodoService struct {
	Repository TodoRepository
}

func (s TodoService) AddTodo(title string, description string) Todo {
	todo, _ := s.Repository.AddTodo(title, description)
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

func (s TodoService) AddAllTodos(requestTodos []Todo) []Todo {
	var todos []Todo
	for _, requestTodo := range requestTodos {
		todo, _ := s.Repository.AddTodo(requestTodo.Title, requestTodo.Description)
		todos = append(todos, todo)
	}
	return todos
}
