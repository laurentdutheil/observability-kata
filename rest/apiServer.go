package rest

import (
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"strconv"
	"todo_odd/repository"
)

type ApiServer struct {
	http.Handler
	repository *repository.TodoRepository
}

var tracer = otel.Tracer("github.com/laurentdutheil/observability-kata/rest")

func NewApiServer() *ApiServer {
	api := &ApiServer{
		repository: &repository.TodoRepository{},
	}

	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", api.HealthcheckHandler)
	router.HandleFunc("/todo/{id}", api.TodoHandlerGet)
	router.HandleFunc("/todo", api.TodoHandler)

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
		_, span := tracer.Start(request.Context(), "todo creation")
		defer span.End()

		var m map[string]interface{}
		_ = json.NewDecoder(request.Body).Decode(&m)
		_ = request.Body.Close()

		todo := s.repository.AddTodo(m["title"].(string), m["description"].(string))
		body, _ := json.Marshal(createJsonTodo(todo))

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
