# 013 - API doc with swaggo

## Status
Proposed

## Context
Consumers of the API need a way to discover available endpoints, request/response shapes and error cases without reading the handler source. Hand-written documentation drifts from the code as endpoints change, and there was no existing mechanism to generate or serve API docs in this project.

## Decision
We use [swaggo](https://github.com/swaggo/swag) to generate an OpenAPI spec from annotations placed directly above each API controller, and serve it through an interactive UI.

## Consequences
Documentation stays close to the code and is easy to regenerate, and the interactive UI lets consumers try requests against the API directly instead of just reading a spec. The trade-off is that every controller needs its own set of annotations kept in sync with the actual handler signature and behavior, which is extra boilerplate and can silently go stale if a handler changes without updating its comments.
