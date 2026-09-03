# Design Patterns

This document lists the design patterns used in the codebase, with a short description of how each is used and pointers to where it is applied.

## Dependency Injection

Constructors accept their collaborators as interfaces or plain arguments instead of constructing them internally, so implementations can be swapped (e.g. for tests) without touching the consumer.

- `catalog.NewCatalogApp(logger, productsRepo, categoriesRepo)` — app/catalog/app.go:10
- Wired up with concrete implementations — cmd/server/main.go

## Repository

Persistence is hidden behind narrow, domain-facing interfaces; callers work with domain types and never see GORM models or SQL.

- Repositories:
    - `ProductsRepository`/`CategoriesRepository` interfaces — app/catalog/interface.go:11-19
    - GORM-backed implementations — infra/models/products_repository.go, infra/models/categories_repository.go
- Monitor:
    - `Monitor`/`Logger` interface used in the whole project, multiple layers - shared/monitor.go
    - Infra implementations - package infra/monitor

## Builder

Configuration objects expose chainable methods that each return the receiver, letting callers assemble a query or a set of validation rules in a single fluent expression.

- `*Query` chainable filters/sorts/pagination — domain/query/query.go:25-67 (`AddFilter`, `AddFilters`, `AddSort`, `AddPagination`)
- `*validator` chainable rule registration — domain/query/validator.go:58-76 (`AllowFilter`, `AllowSort`)

## Strategy

An interface captures an interchangeable algorithm/behavior, and different implementations are selected at construction time.

- `shared.Logger` — real implementation (infra/monitor/logger.go) vs. silent implementation (infra/monitor/noop.go), selected by the caller (e.g. tests use the noop one)
- Per-field value validators passed into `AllowFilter` — domain/query/validator.go:58 (`validateValue func(value any) error`), with reusable strategies `ValidateString`, `ValidatePrice`, `ValidateBool`, `ValidateInt`, `ValidateNumeric` in domain/query/validator.go:148-206

## Null Object

A no-op implementation satisfies the same interface as the real one but performs no observable action, so callers (mainly tests) don't need nil checks or conditional wiring.

- `noopMonitor` / `noopLogger` implement `shared.Monitor` / `shared.Logger` as silent no-ops — infra/monitor/noop.go:16-52

## Decorator

A type wraps another value implementing the same interface, adding behavior transparently while delegating the original calls.

- `statusRecorder` wraps `http.ResponseWriter` to capture the status code written by the handler, for logging — infra/rest/logging_middleware.go:10-18

## Chain of Responsibility / Middleware

HTTP handlers are wrapped in layers, each able to act before/after delegating to the next one, composed in a fixed order around the router.

- `chainMiddlewares` composes `OperationIdMiddleware`, `NewLoggingMiddleware`, `NewTimeoutMiddleware`, `NewCorsMiddleware` around the mux — infra/api/router.go:26-39
- Individual middlewares — infra/rest/operation_id_middleware.go, infra/rest/logging_middleware.go, infra/rest/timeout_middleware.go, infra/rest/cors_middleware.go

## Facade

A small type presents a simplified, unified interface over an underlying dependency, hiding details callers shouldn't depend on directly.

- `monitor` wraps `Logger` behind a single `Monitor` interface — infra/monitor/monitor.go
- `price` package wraps the third-party `decimal` type behind a domain `Price` type, so the dependency isn't imported elsewhere — domain/price/price.go:1-16

## Value Object

Small, immutable types that wrap a primitive value, enforce domain invariants at construction/parsing time, and compare by value rather than identity.

- `ProductCode` — domain/catalog/product_code.go:17-32 (format `PROD###`, validated max value)
- `CategoryCode` — domain/catalog/category_code.go:17-32 (format `CAT###`, validated max value)
- `Price` — domain/price/price.go:14-35 (parses/validates a decimal string)

## Data Transfer Object / Mapper

REST-facing structs are kept separate from domain types, with explicit mapping functions converting between the two at the API boundary.

- `product`, `variant`, `category` REST structs and `encodeProductResponse`/`encodeVariantResponse` mappers — infra/api/product_detail.go:29-68
- `productToDomain`/`productsToDomain` mapping GORM models back to domain types — infra/models/products_repository.go:106,127

## Template Method (via generics + injected functions)

A generic function fixes the shape of an operation (decode → call → encode → write) while the actual decode/call/encode steps are supplied by the caller as functions, avoiding duplication across handlers.

- `NewHandler` — decode request, call app handler, encode response, write it — infra/rest/handler.go:31-55
- `NewListByQueryHandle` — same shape specialized for query-based list endpoints — infra/rest/handler_query.go:14-36

## Read/Write Splitting (CQRS-style connection routing)

Database connections are split into a write source and read replica(s) at the infrastructure layer via GORM's `dbresolver` plugin, so read-heavy queries can be routed to a replica without any change to repository code.

- `database.New` configures a write source and read replica(s) with `dbresolver.RandomPolicy` — infra/database/pg.go:36-82
- Separate `Config` for write/read connection parameters, shared `PoolConfig` for pool tuning — infra/database/pg.go:16-23, infra/database/pool.go:14-19
