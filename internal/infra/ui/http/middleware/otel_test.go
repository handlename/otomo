package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestOtel_Wrap(t *testing.T) {
	// Set up a test TracerProvider
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	// Set up a propagator
	propagator := propagation.TraceContext{}
	otel.SetTextMapPropagator(propagator)

	middleware := NewOtel()

	t.Run("propagates parent trace context from headers", func(t *testing.T) {
		// Generate a valid parent trace span context
		parentTraceID, err := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
		assert.NoError(t, err)
		parentSpanID, err := trace.SpanIDFromHex("0102030405060708")
		assert.NoError(t, err)

		parentSpanCtx := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    parentTraceID,
			SpanID:     parentSpanID,
			TraceFlags: trace.FlagsSampled,
		})

		ctx := trace.ContextWithSpanContext(context.Background(), parentSpanCtx)
		req := httptest.NewRequest(http.MethodPost, "/event", nil)

		// Inject parent context into request headers
		propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

		var innerSpan trace.Span
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			innerSpan = trace.SpanFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		rec := httptest.NewRecorder()
		middleware.Wrap(nextHandler).ServeHTTP(rec, req)

		assert.NotNil(t, innerSpan)
		sc := innerSpan.SpanContext()
		assert.True(t, sc.IsValid())
		assert.Equal(t, parentTraceID, sc.TraceID(), "should propagate trace ID from headers")
		assert.NotEqual(t, parentSpanID, sc.SpanID(), "should generate a new child span ID, not reuse parent span ID")
	})
}
