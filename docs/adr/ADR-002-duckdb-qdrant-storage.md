# ADR-002: Metadata and Retrieval Storage

- Status: Accepted
- Date: 2026-05-19
- Context: Need analytics-friendly metadata storage and semantic retrieval.
- Decision: Plan DuckDB for evidence/audit analytics and Qdrant for semantic/vector search.
- Consequences: Two datastores with clear separation of concerns.
- Alternatives Considered: Postgres-only, Elasticsearch-only.
