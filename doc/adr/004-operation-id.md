# 004 - Operation ID

## Status
Proposed

## Context
When something goes wrong, we need a way to tie a client-reported error back to the exact request that caused it, and to correlate the logs it produced across the stack. A per-request identifier generated ad-hoc in each handler would be inconsistent and easy to forget, and REST-specific solutions (e.g. reading it off the HTTP request) wouldn't cover other entrypoints like a CLI or an event consumer.

## Decision
A dedicated middleware generates a UUIDv7 operation ID and stores it on the request context for every operation, regardless of entrypoint (REST today, with CLI and event entrypoints able to reuse the same mechanism). The ID is returned in the `metadata` of every response and attached to every log line written within that context, so it doubles as a correlation ID across logs and client-facing error reports.

```json
{
  "metadata": {
    "operationId": "b3c1e2a4-1234-4a8b-9e21-abcdef123456"
    // ...
  },
  "data": {
    // ...
  }
}
```

## Consequences
Clients can quote the operation ID when reporting an error, and we can grep logs for that ID to see everything that happened during that request, which speeds up debugging and error tracking. Generating it per-entrypoint keeps the mechanism transport-agnostic, at the cost of each entrypoint needing to wire up its own middleware/context propagation instead of getting it for free. It also relies on every log call going through `WithContext` to pick up the ID — a log line written without the context silently loses the correlation.
