package rest

import (
	"encoding/json"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"strconv"
	"todo_odd/domain"
	"todo_odd/repository"
)

type ApiServer struct {
	http.Handler
	service domain.TodoService
}

func NewApiServer() *ApiServer {

	todoRepository := &repository.TodoRepository{}
	api := &ApiServer{
		service: domain.TodoService{Repository: todoRepository},
	}

	router := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		router.Handle(pattern, handler)
	}

	handleFunc("/healthcheck", api.HealthcheckHandler)
	handleFunc("/todo/{id}", api.TodoHandlerGet)
	handleFunc("/todo", api.TodoHandlerAdd)
	handleFunc("GET /todo-list", api.TodoHandlerGetAll)
	handleFunc("POST /todo-list", api.TodoHandlerAddAll)

	api.Handler = otelhttp.NewHandler(router, "/")
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

	todo := s.service.AddTodo(request.Context(), bodyTodo.Title, bodyTodo.Description)

	writer.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(createJsonTodo(todo))
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerGetAll(writer http.ResponseWriter, _ *http.Request) {
	todos := s.service.GetAll()

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

	todos := s.service.AddAllTodos(request.Context(), requestTodos)

	var responseTodos []Todo
	for _, todo := range todos {
		responseTodos = append(responseTodos, createJsonTodo(todo))
	}

	writer.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(responseTodos)
	_, _ = writer.Write(body)
}
