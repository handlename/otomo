package trace

import (
	"context"

	"github.com/handlename/otomo/config"
	"github.com/morikuni/failure/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes OpenTelemetry tracing based on config.
func InitTracer(ctx context.Context) (func(context.Context) error, error) {
	if !config.Config.Otel.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	var exp sdktrace.SpanExporter
	var err error
	switch config.Config.Otel.GetExporter() {
	case "stdout":
		exp, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to create stdout trace exporter"))
		}
	case "otlp":
		exp, err = otlptracehttp.New(ctx)
		if err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to create otlp trace exporter"))
		}
	default:
		return nil, failure.New(
			failure.Messagef("unsupported otel exporter: %s", config.Config.Otel.GetExporter()),
		)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.Config.Otel.GetServiceName()),
		),
	)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create otel resource"))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}

// RecordError records the error and sets span status to error if err is not nil.
func RecordError(span trace.Span, err error) {
	if err != nil && span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
