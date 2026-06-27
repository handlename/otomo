package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var _ Middleware = (*Otel)(nil)

type Otel struct {
	tracer trace.Tracer
}

func NewOtel() *Otel {
	return &Otel{
		tracer: otel.Tracer("otomo"),
	}
}

// Wrap implements Middleware.
func (o *Otel) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(req.Context(), propagation.HeaderCarrier(req.Header))
		ctx, span := o.tracer.Start(
			ctx,
			req.Method+" "+req.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
