package repository

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoRepositoryInstrumentation struct {
	ctx    context.Context
	tracer trace.Tracer
	span   trace.Span
}

func startInstrumentation(ctx context.Context, name string) *TodoRepositoryInstrumentation {
	i := &TodoRepositoryInstrumentation{ctx: ctx, tracer: otel.Tracer("")}
	i.ctx, i.span = i.tracer.Start(i.ctx, name)
	return i
}

func (i *TodoRepositoryInstrumentation) stopInstrumentation() {
	i.span.End()
}

func (i *TodoRepositoryInstrumentation) todoCreated(id int) {
	i.span.SetAttributes(attribute.Int("id", id))
}
