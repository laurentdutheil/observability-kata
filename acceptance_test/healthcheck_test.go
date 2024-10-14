package acceptance_test

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_odd/rest"
)

func TestHappyHealthcheck(t *testing.T) {
	server := rest.NewApiServer(nil)

	request, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	b, _ := io.ReadAll(response.Body)
	assert.JSONEq(t, `{"status": "OK", "messages": []}`, string(b))
}
