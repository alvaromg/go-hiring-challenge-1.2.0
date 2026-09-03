# 011 - API routes with version prefix

## Status
Proposed

## Context
The API will inevitably need breaking changes to its routes or payloads as it evolves, but existing clients can't be forced to upgrade the moment that happens. Without a versioning scheme, a breaking change would either force all clients to migrate at once or require ad-hoc, per-endpoint workarounds to keep old behavior alive.

## Decision
Every route is registered under a `/v1` prefix (e.g. `/v1/catalog`, `/v1/categories`), as defined in `infra/api/router.go`.

## Consequences
It makes it easy to add new endpoint versions in the future while clients are migrating, since a `/v2` route can be introduced alongside `/v1` without touching or breaking the existing one. The trade-off is that every route carries the prefix from day one, and nothing yet enforces how a version is retired or how long two versions must be supported side by side.
