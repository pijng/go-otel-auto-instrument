package gootelinstrument

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	err := setupTracerProvider(serviceName)
	if err != nil {
		panic(err)
	}
}

func SetContext(ctx context.Context) {
	setContext(ctx)
}

func StartSpan(name string) (context.Context, context.Context, trace.Span) {
	ctx := getContext()
	parentCtx := getParentContext()

	return startSpan(ctx, parentCtx, name)
}

func StartSpanCtx(ctx context.Context, name string) (context.Context, context.Context, trace.Span) {
	parentCtx := getParentContext()
	return startSpan(ctx, parentCtx, name)
}

func startSpan(ctx context.Context, _ context.Context, name string) (context.Context, context.Context, trace.Span) {
	t := otel.Tracer(name)
	gCtx := getContext()

	ctxSpan := trace.SpanFromContext(ctx)
	gCtxSpan := trace.SpanFromContext(gCtx)

	if ctxSpan != nil && ctxSpan.SpanContext().IsValid() {
		newCtx, span := t.Start(ctx, name)
		setContext(newCtx)
		return newCtx, ctx, span
	}

	if gCtxSpan != nil && gCtxSpan.SpanContext().IsValid() {
		newCtx, span := t.Start(gCtx, name)
		setContext(newCtx)
		return newCtx, gCtx, span
	}

	newCtx, span := t.Start(ctx, name)
	setContext(newCtx)
	return newCtx, ctx, span
}
