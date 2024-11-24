package gootelinstrument

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	err := setupTracerProvider("localhost:4318")
	if err != nil {
		panic(err)
	}
}

func StartSpan(name string) trace.Span {
	ctx := getContext()
	parentCtx := getParentContext()

	return startSpan(ctx, parentCtx, name)
}

func StartSpanCtx(ctx context.Context, name string) trace.Span {
	parentCtx := getParentContext()
	return startSpan(ctx, parentCtx, name)
}

func startSpan(ctx context.Context, _ context.Context, name string) trace.Span {
	t := otel.Tracer(name)
	gCtx := getContext()

	ctxSpan := trace.SpanFromContext(ctx)
	gCtxSpan := trace.SpanFromContext(gCtx)

	if ctxSpan != nil && ctxSpan.SpanContext().IsValid() {
		newCtx, span := t.Start(ctx, name)
		setContext(newCtx)
		return span
	}

	if gCtxSpan != nil && gCtxSpan.SpanContext().IsValid() {
		newCtx, span := t.Start(gCtx, name)
		setContext(newCtx)
		return span
	}

	newCtx, span := t.Start(ctx, name)
	setContext(newCtx)
	return span
}
