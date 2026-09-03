# 010 - Create multiple categories at once

## Status
Proposed

## Context
Clients often need to seed several categories together (e.g. when bootstrapping a catalog), and calling a single-category create endpoint once per item is chatty and gives no way to guarantee the whole batch lands together. We needed a way to create several categories in one request without leaving the database in a half-created state if one of them fails.

## Decision
The create category endpoint accepts a list of categories in a single request. All records are inserted inside one database transaction, so if any record fails to be created, none of them are persisted.

## Consequences
Clients get an atomic, all-or-nothing batch create and avoid the overhead of one round-trip per category. The trade-off is that a single invalid record in a large batch rejects the whole batch, and the transaction holds locks across every row for the duration of the insert, which can become a bottleneck for very large batches.
