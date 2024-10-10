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
