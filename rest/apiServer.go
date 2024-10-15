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
	handleFunc("/todo", api.TodoHandler)

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

func (s ApiServer) TodoHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		todos := s.repository.All()
		var bodyResponse []Todo
		for _, todo := range todos {
			bodyResponse = append(bodyResponse, createJsonTodo(todo))
		}
		body, _ := json.Marshal(bodyResponse)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(body)
	case "POST":
		var m map[string]interface{}
		_ = json.NewDecoder(request.Body).Decode(&m)
		_ = request.Body.Close()

		todo := s.repository.AddTodo(m["title"].(string), m["description"].(string))
		body, _ := json.Marshal(createJsonTodo(todo))

		span := trace.SpanFromContext(request.Context())
		span.SetName("todo creation")
		span.SetAttributes(attribute.Int("id", todo.Id))

		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write(body)

	}
}

func (s ApiServer) TodoHandlerGet(writer http.ResponseWriter, request *http.Request) {
	pathId := request.PathValue("id")
	todoId, _ := strconv.Atoi(pathId)

	todo := s.repository.Get(todoId)
	body, _ := json.Marshal(createJsonTodo(todo))
	_, _ = writer.Write(body)

	writer.WriteHeader(http.StatusOK)
}
