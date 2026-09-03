# 009 - Products-categories as many-to-one

## Status
Proposed

## Context
Products need to be organized by category, and we had to decide how tightly to couple the two: a many-to-one relationship (each product belongs to a single category) or a many-to-many relationship (a product can belong to several categories). The many-to-many option is more flexible but adds a join table and more complex queries for a requirement that, so far, only calls for one category per product.

## Decision
We model the products-categories relationship as many-to-one: each product has exactly one category.

## Consequences
This keeps the schema and queries simple and slightly improves database performance compared to a many-to-many approach. The trade-off is that a product can't belong to more than one category; if that becomes a requirement, it will need a schema migration and code changes to move to a many-to-many relationship.
