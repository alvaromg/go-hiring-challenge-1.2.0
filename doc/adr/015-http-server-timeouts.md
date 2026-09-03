# 015 - HTTP server timeouts

## Status
Proposed

## Context
The API server was running with Go's zero-value `http.Server` timeouts, meaning a slow or stalled client (or a handler stuck on a downstream call) could hold a connection open indefinitely, tying up server resources and leaving no bound on request latency. We needed a way to cap both how long the server waits on a client and how long a handler is allowed to run, without hardcoding values that differ between environments.

## Decision
We added `ReadHeaderTimeout` and `IdleTimeout` to both `http.Server`s in `cmd/server/main.go`, and a new `rest.NewTimeoutMiddleware` (`infra/rest/timeout_middleware.go`) that wraps the router with `http.TimeoutHandler` to bound total handler execution time. All three are configurable via `HTTP_READ_HEADER_TIMEOUT`, `HTTP_IDLE_TIMEOUT`, and `HTTP_HANDLER_TIMEOUT` env vars, each falling back to a sane default (5s, 60s, 10s respectively) when unset.

## Consequences
The server now has a predictable upper bound on connection and request lifetime, protecting it from slow-client and hung-handler resource exhaustion, and the handler timeout cancels the request context so downstream calls (e.g. DB queries) that respect `ctx` are cancelled too. The trade-off is that a genuinely slow but valid request (a large report, a slow third-party call) now gets cut off at the configured timeout with a 503 instead of eventually completing, so the values need tuning per endpoint/environment rather than being a one-size-fits-all setting — today it's a single global timeout for every route.
