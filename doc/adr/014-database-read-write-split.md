# 014 - Database read/write split

## Status
Proposed

## Context
As read traffic grows, a single primary Postgres instance becomes a bottleneck and a scaling constraint, while adding a read replica is a standard way to offload reads without touching the primary. The application needed a way to route writes to the primary and reads to a replica without spreading connection/routing logic across every query call site.

## Decision
We split the database connection into separate write and read `Config`s (`infra/database/pg.go`) and register both with gorm's [`dbresolver`](https://github.com/go-gorm/dbresolver) plugin, which transparently sends writes to the source and reads to the replica on the same `*gorm.DB` handle.

## Consequences
Application and repository code keeps using a single `*gorm.DB` without knowing about the split, and adding a replica for read scaling requires no query-level changes. The trade-off is the operational cost of running and keeping a replica in sync, the risk of replication lag surfacing as stale reads right after a write (e.g. read-your-writes scenarios), and one more moving part to configure and monitor. Today write and read configs point at the same instance in every environment (`.env`, `test/helper/api.go`), so the split has no effect until an actual replica exists.
