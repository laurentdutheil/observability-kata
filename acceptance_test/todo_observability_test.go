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
	span := getFirstSpanWithName(t, spans, "todo creation")
	todoResponse := parseTodoResponse(response)
	hasAttribute(t, span, "id", attribute.IntValue(todoResponse.Id))
}

func TestTraceTodoCreationAllNestedSpans(t *testing.T) {
	inMemoryExporter := tracetest.NewInMemoryExporter()
	provider := trace.NewTracerProvider(trace.WithSyncer(inMemoryExporter))
	otel.SetTracerProvider(provider)

	server := rest.NewApiServer()

	todosForPost := validSeveralTodosForPost()
	request, _ := http.NewRequest(http.MethodPost, "/todo-list", bytes.NewBuffer(todosForPost))
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	spans := inMemoryExporter.GetSpans()
	span := getFirstSpanWithName(t, spans, "todo creation all")
	parallelSpans := getSpansWithName(t, spans, "todo creation repo")
	assert.NotEmpty(t, parallelSpans)
	for _, parallelSpan := range parallelSpans {
		assert.Equal(t, span.SpanContext.SpanID(), parallelSpan.Parent.SpanID())
	}
}

func getFirstSpanWithName(t *testing.T, spans tracetest.SpanStubs, spanName string) tracetest.SpanStub {
	foundSpans := getSpansWithName(t, spans, spanName)
	if len(foundSpans) > 0 {
		return foundSpans[0]
	}
	assert.Fail(t, fmt.Sprintf("No span with name '%s' found", spanName))
	return tracetest.SpanStub{}
}

func getSpansWithName(t *testing.T, spans tracetest.SpanStubs, spanName string) (foundSpans []tracetest.SpanStub) {
	require.NotEmpty(t, spans)
	for _, span := range spans {
		if spanName == span.Name {
			foundSpans = append(foundSpans, span)
		}
	}
	return
}

func hasAttribute(t *testing.T, span tracetest.SpanStub, key attribute.Key, value attribute.Value) bool {
	for _, keyValue := range span.Attributes {
		if keyValue.Key == key {
			return assert.Equal(t, keyValue.Value, value)
		}
	}

	return assert.Fail(t, fmt.Sprintf("No attribute with key '%s' found", key))
}
