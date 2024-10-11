package acceptance_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_odd/rest"
)

func TestAddValidTodo(t *testing.T) {
	bodyPost := createValidTodo()
	server := rest.NewApiServer()

	response := postTodoCreation(server, bodyPost)

	assert.Equal(t, http.StatusCreated, response.Code)
	m := parseBodyResponse(response)
	assert.NotNil(t, m["id"])
	assert.Equal(t, "New Todo", m["title"])
	assert.Equal(t, "Description of the todo", m["description"])
}

func parseBodyResponse(response *httptest.ResponseRecorder) map[string]interface{} {
	var m map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&m)
	return m
}

func postTodoCreation(server *rest.ApiServer, bodyPost []byte) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer(bodyPost))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func createValidTodo() []byte {
	bodyPost := []byte(`{
		"title": "New Todo",
		"description": "Description of the todo"
	}`)
	return bodyPost
}