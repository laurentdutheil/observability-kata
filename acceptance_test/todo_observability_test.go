package acceptance_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	spans := inMemoryExporter.GetSpans()
	span := hasSpanWithName(t, spans, "todo creation")
	todoResponse := parseTodoResponse(response)
	hasAttribute(t, span, "id", attribute.IntValue(todoResponse.Id))
}

func hasSpanWithName(t *testing.T, spans tracetest.SpanStubs, spanName string) tracetest.SpanStub {
	require.NotEmpty(t, spans)
	for _, span := range spans {
		if assert.Equal(t, spanName, span.Name) {
			return span
		}
	}
	assert.Fail(t, fmt.Sprintf("No span with name '%s' found", spanName))
	return tracetest.SpanStub{}
}

func hasAttribute(t *testing.T, span tracetest.SpanStub, key attribute.Key, value attribute.Value) bool {
	for _, keyValue := range span.Attributes {
		if keyValue.Key == key {
			return assert.Equal(t, keyValue.Value, value)
		}
	}

	return assert.Fail(t, fmt.Sprintf("No attribute with key '%s' found", key))
}
