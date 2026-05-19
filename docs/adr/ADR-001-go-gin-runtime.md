# ADR-001: Go Runtime and HTTP Service Choice

- Status: Accepted
- Date: 2026-05-19
- Context: Need a robust backend for evidence-chain APIs and security-first testing.
- Decision: Use Go as runtime; current implementation uses standard `net/http` for maximal portability in this environment.
- Consequences: Strong performance and simple deployment; framework features can be added later with minimal API-surface changes.
- Alternatives Considered: Gin, Fiber, Node/NestJS.
