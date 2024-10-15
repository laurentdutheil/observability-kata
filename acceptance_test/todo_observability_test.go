package acceptance_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_odd/rest"
)

func TestTraceTodoCreation(t *testing.T) {
	inMemoryExporter := tracetest.NewInMemoryExporter()
	provider := trace.NewTracerProvider(trace.WithSyncer(inMemoryExporter))
	otel.SetTracerProvider(provider)

	server := rest.NewApiServer()

	bodyPost := validTodoForPost()
	request, _ := http.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer(bodyPost))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	traces := inMemoryExporter.GetSpans()
	assert.NotEmpty(t, traces)
	assert.Equal(t, "todo creation", traces[0].Name)
	spanAttribute := traces[0].Attributes[0]
	assert.Equal(t, attribute.Key("id"), spanAttribute.Key)
	todoResponse := parseTodoResponse(response)
	assert.Equal(t, spanAttribute.Value.AsInt64(), int64(todoResponse.Id))
}
