# Tools Integration Blueprint

This project integrates and/or stages integration for the curated stack sourced from `/Users/dennis_leedennis_lee/Desktop/software-tools`:

- MCP + OpenAPI Tool Calling: planned via `openapi/openapi.yaml` and MCP adapter backlog.
- Ollama: planned as local AI enrichment provider for evidence classification/summarization.
- Qdrant: planned as vector retrieval backend for semantic evidence search.
- DuckDB: planned as local analytics/audit metadata store.
- Weights & Biases Weave: planned as AI evaluation + observability layer for enrichment quality.
- OpenTelemetry: planned for traces/metrics across API, storage, and AI adapters.
- Cilium Tetragon: planned runtime security detection profiles and policy mapping.
- PQC libraries: planned post-quantum signing placeholders and dual-sign strategy.

Current implementation baseline in this commit:
- Go/Gin API platform scaffold.
- OpenAPI spec for core APIs.
- Robot Framework suites including OWASP Top10 baseline tests.
- CI/CD pipelines for Go tests, OpenAPI parse validation, Robot/OWASP execution.
