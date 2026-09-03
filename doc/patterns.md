# Design Patterns

Dependency Injection — constructors accept their collaborators as interfaces/args (e.g. `catalog.NewCatalogApp(logger, productsRepo, categoriesRepo)` in app/catalog/app.go, wired up in cmd/server/main.go)
Repository — `ProductsRepository`/`CategoriesRepository` interfaces in app/catalog/interface.go, implemented by GORM-backed structs in infra/models
Builder — fluent, chainable configuration objects returning `*Query`/`*validator` (domain/query/query.go, domain/query/validator.go)
Strategy — pluggable `shared.Logger` implementations (infra/monitor/logger.go vs infra/monitor/noop.go) and pluggable filter-value validators passed into `AllowFilter` (domain/query/validator.go)
Null Object — `noopMonitor`/`noopLogger` implement the real interfaces as silent no-ops for tests (infra/monitor/noop.go)
Decorator — `statusRecorder` wraps `http.ResponseWriter` to capture the status code (infra/rest/logging_middleware.go)
Chain of Responsibility / Middleware — `chainMiddlewares` composes `OperationIdMiddleware`, `NewLoggingMiddleware`, `NewCorsMiddleware` around the router (infra/api/router.go)
Facade — `Monitor` wraps `Logger` behind a small unified interface (infra/monitor/monitor.go)
Value Object — immutable domain identifiers/amounts with parsing and invariants, e.g. `ProductCode`, `CategoryCode`, `Price` (domain/catalog/product_code.go, domain/price/price.go)
Data Transfer Object / Mapper — REST-facing structs (`product`, `variant`, `category`) mapped to/from domain types via `encode*Response`/`*ToDomain` functions (infra/api/product_detail.go, infra/models/products_repository.go)
Template Method (via generics + injected functions) — `NewHandler`/`NewListByQueryHandle` define a fixed decode → call → encode → write pipeline parameterized by request decoder, app handler and response encoder (infra/rest/handler.go, infra/rest/handler_query.go)
DB reads/writes split
