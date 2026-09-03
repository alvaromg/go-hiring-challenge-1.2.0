# 003 - Generic list

## Status
Proposed

## Context
Several resources (products, categories, ...) need list endpoints that return multiple items alongside pagination metadata such as total count and page info. Building this per resource means re-deriving the same envelope shape and page-count math each time, with room for inconsistencies between endpoints in how totals and pages are reported. We needed a single reusable representation for "a page of items" that any resource could plug into.

## Decision
We use a generic `list.ListResponse[T]` in `domain/list` to carry items and their total count through the application layer, and a generic `EncodeListResponse` in `infra/rest` to turn it into the HTTP response, wrapping items with `[]RO -> WD` so each resource can nest them under its own key (e.g. `products`). The response metadata carries `page`, `pageSize`, `totalCount` and a derived `pageCount`, alongside the standard `operationId` from [[002-api-responses-format]].

```json
{
  "metadata": {
    "operationId": "b3c1e2a4-1234-4a8b-9e21-abcdef123456",
    "totalCount": 57,
    "page": 2,
    "pageSize": 20,
    "pageCount": 3
  },
  "data": {
    "products": [
      { "code": "PROD123", "price": "129.90" }
    ]
  }
}
```

## Consequences
Every list endpoint returns pagination metadata in the same shape, computed in one place, so resources can't drift on how totals or page counts are reported. The generic type parameters (`DO`, `RO`, `WD`) add the same indirection and compiler-error verbosity noted in [001-generic-api-handlers.md](001-generic-api-handlers.md), and any listing that doesn't fit the "items + total" shape (e.g. cursor-based or streaming pagination) would need a separate path outside this abstraction.
