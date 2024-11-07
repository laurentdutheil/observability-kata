package domain

type TodoRepository interface {
	AddTodo(title string, description string) (Todo, error)
	Get(id int) (Todo, error)
	All() ([]Todo, error)
}
