package rest

import (
	"encoding/json"
	"net/http"
)

type ApiServer struct {
	http.Handler
}

func NewApiServer() *ApiServer {
	server := &ApiServer{}

	router := http.NewServeMux()

	router.HandleFunc("/healthcheck", HealthcheckHandler)
	router.HandleFunc("/todo", TodoHandler)

	server.Handler = router
	return server
}

func HealthcheckHandler(writer http.ResponseWriter, _ *http.Request) {
	health := Health{
		Status:   "OK",
		Messages: []string{},
	}

	body, _ := json.Marshal(health)

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(body)
}

func TodoHandler(writer http.ResponseWriter, request *http.Request) {
	var m map[string]interface{}
	_ = json.NewDecoder(request.Body).Decode(&m)
	_ = request.Body.Close()

	todo := Todo{
		Id:          1,
		Title:       m["title"].(string),
		Description: m["description"].(string),
	}
	body, _ := json.Marshal(todo)

	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(body)
}
