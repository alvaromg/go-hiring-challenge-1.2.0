# 016 - Observability with OpenTelemetry and Grafana

## Status
Proposed

## Context
Structured logs ([005-structured-logs](005-structured-logs.md)) tell us what happened at a single point, but not how a request's time was spent across the HTTP handler, GORM queries and downstream calls, or how to correlate a log line with the request that produced it. We needed traces and logs to land in one place that can be queried and visualized, without hand-rolling a metrics/tracing backend.

## Decision
We instrument the app with OpenTelemetry: `monitor.SetupOTelSDK` (`infra/monitor/otel.go`) bootstraps an OTLP/HTTP trace and log pipeline, `infra/rest/tracing.go` starts a server span per request, and `infra/database/tracing.go` wraps every GORM statement in a client span, so DB spans nest under the HTTP span that triggered them. Logs are additionally piped through the same OTel logger provider via an `otellogrus` hook. Everything is exported to a `grafana-lgtm` container (Grafana + Loki + Tempo + Mimir bundled in one image) added to `docker-compose.yml`, configured via `OTEL_EXPORTER_OTLP_ENDPOINT`.

## Consequences
We get traces and logs correlated in one Grafana instance (`http://localhost:3000`) with a single `docker-compose up`, and no separate backend to provision. The trade-offs: `grafana-lgtm` is a local/dev convenience, not something we'd run as-is in production; every DB call now carries the overhead of a span even when nothing is exporting; and the tracer is a required constructor argument threaded through `monitor.NewMonitor` and `database.New`, so every caller (including tests, via `monitor.NewNoopMonitor`) has to supply one, even a no-op.
