package monitor

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	otellog "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

// SetupOTelSDK bootstraps the OpenTelemetry trace and log pipelines,
// exporting both via OTLP/HTTP to otlpEndpoint (host:port, no scheme -
// e.g. "localhost:4318" for the grafana-lgtm container). It registers the
// resulting tracer provider and logger provider as the global providers, so
// packages that only have access to the global otel API (such as instrumented
// libraries) pick them up automatically.
//
// The returned shutdown function flushes and closes both exporters and
// should be called once during graceful shutdown.
func SetupOTelSDK(ctx context.Context, serviceName, otlpEndpoint string) (tracer trace.Tracer, loggerProvider *sdklog.LoggerProvider, shutdown func(context.Context) error, err error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(otlpEndpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
	otellog.SetLoggerProvider(lp)

	shutdown = func(ctx context.Context) error {
		return errors.Join(
			tracerProvider.Shutdown(ctx),
			lp.Shutdown(ctx),
		)
	}

	return tracerProvider.Tracer(serviceName), lp, shutdown, nil
}
