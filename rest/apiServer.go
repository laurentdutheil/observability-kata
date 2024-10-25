package rest

import (
	"encoding/json"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"todo_odd/repository"
)

type ApiServer struct {
	http.Handler
	repository *repository.TodoRepository
}

func NewApiServer() *ApiServer {
	api := &ApiServer{
		repository: &repository.TodoRepository{},
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

func (s ApiServer) TodoHandlerAdd(writer http.ResponseWriter, request *http.Request) {
	var bodyTodo Todo
	_ = json.NewDecoder(request.Body).Decode(&bodyTodo)
	_ = request.Body.Close()

	todo := s.repository.AddTodo(request.Context(), bodyTodo.Title, bodyTodo.Description)
	body, _ := json.Marshal(createJsonTodo(todo))

	span := trace.SpanFromContext(request.Context())
	span.SetName("todo creation")
	span.SetAttributes(attribute.Int("id", todo.Id))

	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerGetAll(writer http.ResponseWriter, _ *http.Request) {
	todos := s.repository.All()
	var bodyResponse []Todo
	for _, todo := range todos {
		bodyResponse = append(bodyResponse, createJsonTodo(todo))
	}
	body, _ := json.Marshal(bodyResponse)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerAddAll(writer http.ResponseWriter, request *http.Request) {
	var bodyTodos []Todo
	_ = json.NewDecoder(request.Body).Decode(&bodyTodos)
	_ = request.Body.Close()

	span := trace.SpanFromContext(request.Context())
	span.SetName("todo creation all")

	var todos []Todo
	for _, bodyTodo := range bodyTodos {
		todo := s.repository.AddTodo(request.Context(), bodyTodo.Title, bodyTodo.Description)
		todos = append(todos, createJsonTodo(todo))
	}

	body, _ := json.Marshal(todos)

	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(body)
}

func (s ApiServer) TodoHandlerGet(writer http.ResponseWriter, request *http.Request) {
	pathId := request.PathValue("id")
	todoId, _ := strconv.Atoi(pathId)

	todo := s.repository.Get(todoId)
	body, _ := json.Marshal(createJsonTodo(todo))
	_, _ = writer.Write(body)

	writer.WriteHeader(http.StatusOK)
}
