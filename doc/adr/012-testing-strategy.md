# 012 - Testing strategy

## Status
Proposed

## Context
The codebase mixes domain logic with layers that talk to real infrastructure (API, application services, repositories), and each needs a different testing approach to give confidence without becoming slow or brittle. The main choice was how much to mock: heavy mocking keeps tests fast but couples them to implementation details and adds boilerplate, while testing against real dependencies is closer to production behavior at the cost of speed.

## Decision
Unit tests cover domain logic that has no external dependencies. Integration tests exercise the full stack from the API down through the application layer and repositories, without mocking, using gnomock containers with the same data schema as production. Integration tests live in `/test` and are grouped by similarity so related tests can share a container instead of spinning up one per test.

## Consequences
This avoids the boilerplate that mocking every dependency tends to add, and integration tests give real confidence that the layers are wired together correctly, closer to actual production behavior. The trade-off is slower test runs than a fully mocked suite, a dependency on Docker for gnomock containers, and grouping tests to share containers can make failures harder to isolate since more code runs before a single assertion fails.
