package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"todo_odd/domain"
)

type constructor func() domain.TodoRepository

var constructors = []constructor{
	func() domain.TodoRepository {
		return &InMemoryRepository{}
	},
	func() domain.TodoRepository {
		return NewSqliteRepository()
	},
}

func TestTodoRepository(t *testing.T) {
	for _, newTodoRepository := range constructors {

		t.Run("add todo", func(t *testing.T) {
			repo := newTodoRepository()

			todo, err := repo.AddTodo("Title", "Description")

			assert.NoError(t, err)
			assert.Greater(t, todo.Id, 0)
		})

		t.Run("get todo", func(t *testing.T) {
			repo := newTodoRepository()
			expectedTodo, _ := repo.AddTodo("Title", "Description")

			got, err := repo.Get(expectedTodo.Id)
			assert.NoError(t, err)
			assert.Equal(t, expectedTodo, got)

		})

		t.Run("get inexistant todo", func(t *testing.T) {
			repo := newTodoRepository()

			_, err := repo.Get(0)
			assert.ErrorContains(t, err, "repository: todo #0 does not exist")
		})

		t.Run("get all todos", func(t *testing.T) {
			repo := newTodoRepository()
			for i := 1; i <= 3; i++ {
				_, _ = repo.AddTodo(fmt.Sprintf("Title_%d", i), fmt.Sprintf("Description_%d", i))
			}

			todos, err := repo.All()
			assert.NoError(t, err)
			assert.Len(t, todos, 3)
			for i := 1; i <= 3; i++ {
				assert.Equal(t, fmt.Sprintf("Title_%d", i), todos[i-1].Title)
				assert.Equal(t, fmt.Sprintf("Description_%d", i), todos[i-1].Description)
			}
		})
	}
}
