package rest

import (
	"todo_odd/domain"
)

type Health struct {
	Status   string   `json:"status"`
	Messages []string `json:"messages"`
}

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Error struct {
	Message string `json:"message"`
}

func (t Todo) ToDomainTodo() domain.Todo {
	return domain.Todo{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
	}
}

func createJsonTodo(todo domain.Todo) Todo {
	return Todo{
		Id:          todo.Id,
		Title:       todo.Title,
		Description: todo.Description,
	}
}
