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
	bodyPost := []byte(`{
		"title": "New Todo",
		"description": "Description of the todo"
	}`)

	request, _ := http.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer(bodyPost))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()

	server := rest.NewApiServer()
	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&m)
	assert.NotNil(t, m["id"])
	assert.Equal(t, "New Todo", m["title"])
	assert.Equal(t, "Description of the todo", m["description"])
}
