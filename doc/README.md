# Documentation

## Summary

Go Hiring Challenge is a REST API for browsing a product catalog: products
with prices and variants, organized into categories. It's built as a small
Go service backed by Postgres, with a layered (domain / application / infra)
codebase, database read/write splitting, and Swagger-generated API docs.

## Features

- **Browse the catalog** — list products with offset pagination (page size
  configurable, capped at 100 per page) and the total count of matching
  products.
- **Filter products** — narrow the catalog by category or by a maximum
  price.
- **Sort products** — order results by code or price.
- **Product details** — look up a single product by its code, including its
  category and its variants; variants without their own price inherit the
  product's price.
- **Browse categories** — list all available categories.
- **Create categories** — add one or more new categories in a single
  request.
- **API documentation** — interactive Swagger UI describing every endpoint,
  served alongside the API.

## Further reading

- [Architecture](architecture.md) — how the codebase is layered and how
  dependencies flow between layers.
- [Architecture Decision Records](adr/README.md) — the decisions made along
  the way, the problem each one solved, and the trade-offs involved.
- [Design Patterns](patterns.md) — patterns used across the codebase, with
  pointers to where each is applied.
