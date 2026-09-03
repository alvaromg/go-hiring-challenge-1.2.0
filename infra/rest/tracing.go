package rest

import (
	"context"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// startHandlerSpan starts a server span for an incoming REST request. The
// span is named after the route pattern matched by http.ServeMux (e.g.
// "GET /v1/catalog/{code}"), which keeps the number of distinct span names
// bounded regardless of the path parameters of any single request.
func startHandlerSpan(monitor shared.Monitor, r *http.Request) (context.Context, trace.Span) {
	spanName := r.Pattern
	if spanName == "" {
		spanName = r.Method + " " + r.URL.Path
	}

	ctx, span := monitor.Tracer().Start(r.Context(), spanName, trace.WithSpanKind(trace.SpanKindServer))
	span.SetAttributes(
		attribute.String("http.request.method", r.Method),
		attribute.String("url.path", r.URL.Path),
	)

	return ctx, span
}

// endHandlerSpan records the outcome of the request on the span and ends it.
// It must be called exactly once per span, typically via defer.
func endHandlerSpan(span trace.Span, statusCode int, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.SetAttributes(attribute.Int("http.response.status_code", statusCode))
	span.End()
}
