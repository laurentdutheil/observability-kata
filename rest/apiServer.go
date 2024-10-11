package rest

import (
	"encoding/json"
	"net/http"
	"todo_odd/repository"
)

type ApiServer struct {
	http.Handler
	repository repository.TodoRepository
}

func NewApiServer() *ApiServer {
	api := &ApiServer{}

	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", api.HealthcheckHandler)
	router.HandleFunc("/todo", api.TodoHandler)

	api.Handler = router
	return api
}

func (s ApiServer) HealthcheckHandler(writer http.ResponseWriter, _ *http.Request) {
	health := repository.Health{
		Status:   "OK",
		Messages: []string{},
	}

	body, _ := json.Marshal(health)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandler(writer http.ResponseWriter, request *http.Request) {
	var m map[string]interface{}
	_ = json.NewDecoder(request.Body).Decode(&m)
	_ = request.Body.Close()

	todo := s.repository.AddTodo(m["title"].(string), m["description"].(string))
	body, _ := json.Marshal(todo)

	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(body)
}
