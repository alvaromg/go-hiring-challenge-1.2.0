# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Detailed documentation

- [Documentation index](doc/README.md)
  - [Architecture](doc/architecture.md) — how the codebase is layered and how
    dependencies flow between layers.
  - [Architecture Decision Records](doc/adr/README.md) — the decisions made along
    the way, the problem each one solved, and the trade-offs involved.
  - [Design Patterns](doc/patterns.md) — patterns used across the codebase, with
    pointers to where each is applied.

## Project Structure

The codebase follows a layered architecture — see [Architecture](doc/architecture.md) for how these layers relate.

- **cmd/**: Application entry points.
  - `server/main.go`: Starts the REST API (and, optionally, the Swagger doc server).
  - `seed/main.go`: Seeds the database with initial product and category data.
- **domain/**: Core business entities, value objects and rules (e.g. `catalog.Product`, `price.Price`, `query.Query`), with no dependency on the app or infra layers.
- **app/**: Application layer — use cases that orchestrate domain behavior (`catalog.App`) and the repository interfaces infra implements.
- **infra/**: Concrete implementations that wire the app to the outside world.
  - `infra/api`, `infra/rest`: HTTP routing, handlers and middlewares.
  - `infra/models`, `infra/database`: database repositories and connection setup (incl. read/write split).
  - `infra/monitor`: logging and observability.
- **shared/**: Small cross-cutting interfaces (`Logger`, `Monitor`) used across layers.
- **sql/**: Database migration and seed scripts, applied in order.
- **swagger/**: Generated API documentation (via `make swag`).
- **test/**: Integration tests and test helpers.
- **doc/**: Project documentation — see [Documentation index](doc/README.md).
- `.env`: Environment variables file for configuration.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run all tests.
  - `make test-integration`: Will run integration tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.
  - `make swag`: Will update swagger documentation based in code annotations.

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)

## Local environment URLs

Once the app and its Docker services are running (`make docker-up`, `make run`):

| Service      | URL                     |
| ------------ | ----------------------- |
| API          | http://localhost:8484   |
| Swagger docs | http://localhost:1323/swagger/index.html   |
| Grafana      | http://localhost:3000   |
