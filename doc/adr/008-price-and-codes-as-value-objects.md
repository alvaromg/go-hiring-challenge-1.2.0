# 008 - Price and codes as value objects

## Status
Proposed

## Context
Prices and codes (product and category) carry domain rules — a code has a required format, a price is a decimal amount — that primitive types like `string` or `float64` can't enforce. Using primitives lets invalid values flow through the system unchecked, and requires the same parsing and validation logic to be duplicated at every boundary. We also depend on a third-party decimal library to represent monetary amounts, and importing it directly wherever a price is needed would spread that dependency across the codebase.

## Decision
We model `Price`, `ProductCode` and `CategoryCode` as value objects in the domain layer, each validated and constructed through a single `New`/`Parse` entrypoint. `Price` wraps `decimal.Decimal` behind a domain-owned type in `domain/price`, so `github.com/shopspring/decimal` is only imported there.

## Consequences
Invalid codes and prices are rejected at construction instead of drifting through the system as unchecked primitives, and call sites can trust a `Price` or `ProductCode` value without re-validating it. Swapping the underlying decimal library only touches `domain/price`. The trade-off is more types and boilerplate (constructors, `Parse`, `String`, `MarshalJSON`) compared to using primitives directly, and every layer that handles these values needs to convert to and from them at its boundaries.
