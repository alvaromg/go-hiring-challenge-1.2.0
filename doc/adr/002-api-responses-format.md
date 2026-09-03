# 002 - API responses format

## Status
Proposed

## Context
Endpoints need to return both the actual response payload and information about the request that isn't part of that payload, such as the operation ID used for tracing. Mixing the two into a single flat object forces every consumer to distinguish "real" fields from bookkeeping fields by convention, and makes it awkward to add new bookkeeping fields later without risking a collision with the data itself.

## Decision
Every API response is a generic `Response[T]` envelope with two top-level objects: `metadata`, holding information not directly related to the content (currently the operation ID), and `data`, holding the response payload itself. The envelope is built by a single generic encoder shared by all handlers. Errors follow the same shape, but with an `error` object instead of `data`.

Product detail response:
```json
{
  "metadata": {
    "operationId": "b3c1e2a4-1234-4a8b-9e21-abcdef123456"
  },
  "data": {
    "code": "PROD123",
    "price": "129.90",
    "category": {
      "id": "1",
      "name": "Shoes"
    },
    "variants": [
      { "name": "EU 40", "sku": "PROD123A", "price": "129.90" }
    ]
  }
}
```

Error response:
```json
{
  "metadata": {
    "operationId": "b3c1e2a4-1234-4a8b-9e21-abcdef123456"
  },
  "error": {
    "message": "resource not found"
  }
}
```

## Consequences
Consumers get a stable, predictable shape across all endpoints, and new metadata fields (pagination info, request timing, etc.) can be added without ever risking a collision with response data. The cost is a small amount of nesting and boilerplate on every response, and existing API consumers built around a flat response body would need to change how they read the payload.
