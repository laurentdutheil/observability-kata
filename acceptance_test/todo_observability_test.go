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
	todoResponse := parseTodoResponse(response)
	attributeValue := getAttributeValue(traces[0].Attributes, "id")
	assert.Equal(t, attributeValue.AsInt64(), int64(todoResponse.Id))
}

func getAttributeValue(attributes []attribute.KeyValue, key attribute.Key) *attribute.Value {
	for _, keyValue := range attributes {
		if keyValue.Key == key {
			return &keyValue.Value
		}
	}
	return nil
}
