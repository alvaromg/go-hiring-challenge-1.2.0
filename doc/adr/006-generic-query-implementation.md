# 006 - Generic query implementation

## Status
Proposed

## Context
List endpoints need filters, sorting and pagination, and each resource has different rules for which fields and operators are safe to expose (e.g. filtering products by category but not by internal fields). Handling this per resource means re-parsing query strings, re-building `WHERE`/`ORDER BY` clauses and re-validating allowed fields for every new endpoint, with no shared place to stop an unsafe field or operator from reaching the database.

## Decision
We use a resource-agnostic `query.Query` in `domain/query` to carry filters, sorts and pagination through the application layer. `infra/rest` decodes query string parameters into a `Query`, `infra/database` encodes a `Query` into GORM `Where`/`Order`/`Offset`/`Limit` calls, and a `query.Validator` sits in between, letting each resource explicitly declare which fields and operators it allows for filtering and sorting before the query reaches the database.

`infra/rest.DecodeQueryFromRequest` turns `filter_<field>_<operator>`, `sort` and `page`/`pageSize` query params into a `*query.Query`, so a request like:

```
GET /v1/catalog?filter_category_eq=CAT001&filter_price_lt=100.00&sort=-price&page=1&pageSize=20
```

is decoded into filters `category eq "CAT001"` and `price lt "100.00"`, sort `price desc`, and pagination `page 1 / pageSize 20`, before being passed to the application handler.

In `infra/models/products_repository.go`, `ProductsRepository.GetProducts` builds a validator scoped to what that repository allows before running the query:

```go
validator := query.NewValidator().
	AllowFilter(FieldCategory, []query.Operator{query.Eq}, query.ValidateString).
	AllowFilter(FieldPrice, []query.Operator{query.Lt}, query.ValidatePrice).
	AllowSort("code", "price")
if err := validator.Validate(q); err != nil {
	return zero, err
}
```

Filtering products by `category` is only allowed with `eq`, filtering by `price` only with `lt`, and only `code`/`price` may be sorted on — any other field or operator is rejected before it reaches `ApplyQueryFilters`/`ApplyQuerySorts`.

Because pagination lives on the `Query` itself, list responses get their pagination metadata out of the box: `EncodeListResponse` reads `page`/`pageSize` straight off the `Query` and combines them with the `totalCount` the repository sets via `list.ListResponse.SetTotal` to derive `pageCount`, so no resource has to compute or wire up this metadata by hand (see [003-generic-list](003-generic-list.md)).

## Consequences
Every list endpoint gets filters, sorting and pagination for free with consistent parsing and query-building logic, and the validator gives each resource explicit, opt-in control over what can be filtered or sorted on, closing off unsafe or unintended fields by default. The trade-off is the indirection of passing a `Query` through three layers, and the validator itself is opt-in: a handler that skips wiring one up exposes every field on the underlying model to filtering and sorting. It also only covers the filter/sort/pagination shape implemented here — an endpoint needing full-text search or cursor-based pagination would need a separate path.
