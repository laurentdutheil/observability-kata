package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo_odd/domain"
)

type ApiServer struct {
	http.Handler
	service domain.TodoService
}

func NewApiServer(todoRepository domain.TodoRepository) *ApiServer {
	api := &ApiServer{
		service: domain.TodoService{Repository: todoRepository},
	}

	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", api.HealthcheckHandler)
	router.HandleFunc("/todo/{id}", api.TodoHandlerGet)
	router.HandleFunc("/todo", api.TodoHandlerAdd)
	router.HandleFunc("GET /todo-list", api.TodoHandlerGetAll)
	router.HandleFunc("POST /todo-list", api.TodoHandlerAddAll)

	api.Handler = router
	return api
}

func (s ApiServer) HealthcheckHandler(writer http.ResponseWriter, _ *http.Request) {
	health := Health{
		Status:   "OK",
		Messages: []string{},
	}

	body, _ := json.Marshal(health)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerGet(writer http.ResponseWriter, request *http.Request) {
	pathId := request.PathValue("id")
	todoId, _ := strconv.Atoi(pathId)

	todo, err := s.service.GetTodo(todoId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		body, _ := json.Marshal(Error{Message: err.Error()})
		_, _ = writer.Write(body)
		return
	}

	writer.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(createJsonTodo(todo))
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerAdd(writer http.ResponseWriter, request *http.Request) {
	var bodyTodo Todo
	_ = json.NewDecoder(request.Body).Decode(&bodyTodo)
	_ = request.Body.Close()

	todo := s.service.AddTodo(bodyTodo.Title, bodyTodo.Description)

	writer.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(createJsonTodo(todo))
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerGetAll(writer http.ResponseWriter, _ *http.Request) {
	todos, _ := s.service.GetAll()

	var bodyResponse []Todo
	for _, todo := range todos {
		bodyResponse = append(bodyResponse, createJsonTodo(todo))
	}

	writer.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(bodyResponse)
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerAddAll(writer http.ResponseWriter, request *http.Request) {
	var bodyTodos []Todo
	_ = json.NewDecoder(request.Body).Decode(&bodyTodos)
	_ = request.Body.Close()

	var requestTodos []domain.Todo
	for _, bodyTodo := range bodyTodos {
		requestTodos = append(requestTodos, bodyTodo.ToDomainTodo())
	}

	todos := s.service.AddAllTodos(requestTodos)

	var responseTodos []Todo
	for _, todo := range todos {
		responseTodos = append(responseTodos, createJsonTodo(todo))
	}

	writer.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(responseTodos)
	_, _ = writer.Write(body)
}
