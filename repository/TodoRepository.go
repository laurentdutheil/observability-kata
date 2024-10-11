package repository

type TodoRepository struct {
}

func (r *TodoRepository) AddTodo(title string, description string) Todo {
	todo := Todo{
		Id:          1,
		Title:       title,
		Description: description,
	}
	return todo
}
