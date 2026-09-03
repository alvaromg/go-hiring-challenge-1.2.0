# 007 - Semantic errors

## Status
Proposed

## Context
Domain and application code needs a way to signal *why* an operation failed (validation, not found, duplicated resource) so the infra layer can translate that failure into the right HTTP status code, without infra having to parse error message strings or import domain-specific types. The two common options are sentinel error values checked with `errors.Is`, or custom error types checked with `errors.As`.

## Decision
We define a small set of sentinel errors in `domain/errors` (`ErrorDomainValidation`, `ErrorNotFound`, `ErrorDuplicatedResource`) that domain and infra code wrap with `fmt.Errorf("%w: ...", ...)` to add context. The REST layer maps these to HTTP codes using `errors.Is`.

## Consequences
This keeps error handling simple: adding a new failure case is a one-line wrapped error, and the HTTP mapping stays a short, centralized `errors.Is` chain. The downside is that sentinel errors carry no structured data — if a future error needs to expose extra fields (e.g. which field failed validation, or a list of conflicting IDs) beyond a formatted message, sentinels won't be enough and we'll need custom error implementations (structs implementing `error` and `Unwrap`), checked with `errors.As` instead of `errors.Is`.
