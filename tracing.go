package go_otel_auto_instrument

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func setupTracerProvider(serviceName string) error {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
	)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return fmt.Errorf("creating otel exporter: %w", err)
	}

	provider := sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exporter),
		sdk_trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}
