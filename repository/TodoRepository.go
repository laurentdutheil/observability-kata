package repository

type TodoRepository struct {
	todos []Todo
}

func (r *TodoRepository) AddTodo(title string, description string) Todo {
	todo := Todo{
		Id:          len(r.todos) + 1,
		Title:       title,
		Description: description,
	}
	r.todos = append(r.todos, todo)
	return todo
}

func (r *TodoRepository) Get(id int) Todo {
	return r.todos[id-1]
}
