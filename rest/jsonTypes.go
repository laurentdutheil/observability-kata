package rest

import "todo_odd/repository"

type Health struct {
	Status   string   `json:"status"`
	Messages []string `json:"messages"`
}

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createJsonTodo(rTodo repository.Todo) Todo {
	return Todo{
		Id:          rTodo.Id,
		Title:       rTodo.Title,
		Description: rTodo.Description,
	}
}
