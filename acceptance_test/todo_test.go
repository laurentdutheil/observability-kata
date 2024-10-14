package acceptance_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_odd/rest"
)

func TestAddValidTodo(t *testing.T) {
	server := rest.NewApiServer(nil)

	bodyPost := validTodoForPost()
	request, _ := http.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer(bodyPost))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assertTodoResponse(t, parseTodoResponse(response))
}

func TestGetTodoById(t *testing.T) {
	server := rest.NewApiServer(nil)

	id := createValidTodo(server)

	requestURL := fmt.Sprintf("/todo/%d", id)
	request, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assertTodoResponse(t, parseTodoResponse(response), id)
}

func TestGetTodos(t *testing.T) {
	server := rest.NewApiServer(nil)

	id1 := createValidTodo(server)
	id2 := createValidTodo(server)

	request, _ := http.NewRequest(http.MethodGet, "/todo", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	var todos []rest.Todo
	_ = json.NewDecoder(response.Body).Decode(&todos)
	assert.Len(t, todos, 2)
	assertTodoResponse(t, todos[0], id1)
	assertTodoResponse(t, todos[1], id2)
}

func createValidTodo(server *rest.ApiServer) int {
	bodyPost := validTodoForPost()
	r := postTodoCreation(server, bodyPost)
	createdTodo := parseTodoResponse(r)
	id := createdTodo.Id
	return id
}

func assertTodoResponse(t *testing.T, todoResponse rest.Todo, expectedId ...int) {
	if len(expectedId) == 0 {
		assert.NotNil(t, todoResponse.Id)
	} else {
		assert.Equal(t, expectedId[0], todoResponse.Id)
	}
	assert.Equal(t, "New Todo", todoResponse.Title)
	assert.Equal(t, "Description of the todo", todoResponse.Description)
}

func parseTodoResponse(response *httptest.ResponseRecorder) rest.Todo {
	var todo rest.Todo
	_ = json.NewDecoder(response.Body).Decode(&todo)
	return todo
}

func postTodoCreation(server *rest.ApiServer, bodyPost []byte) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer(bodyPost))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func validTodoForPost() []byte {
	bodyPost := []byte(`{
		"title": "New Todo",
		"description": "Description of the todo"
	}`)
	return bodyPost
}
