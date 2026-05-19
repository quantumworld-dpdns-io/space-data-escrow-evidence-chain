# Dev Bootstrap and CLI

## Bootstrap

```bash
make bootstrap
```

This runs:
- `go mod tidy`
- local Robot Framework tool bootstrap

## CLI Commands

Use the CLI against a running API:

```bash
go run ./cmd/cli evidence-create <external_id> <source> <type> <k=v,k2=v2>
go run ./cmd/cli evidence-verify <evidence_id>
go run ./cmd/cli audit-query
go run ./cmd/cli enrich-trigger <evidence_id>
```

Environment variables:
- `EVIDENCE_API_URL` (default `http://localhost:8080`)
- `EVIDENCE_API_KEY` (default `dev-api-key`)
