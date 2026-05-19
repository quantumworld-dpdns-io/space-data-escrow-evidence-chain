# space-data-escrow-evidence-chain

Space data escrow and evidence chain platform with verifiable chain-of-custody APIs, audit trails, security testing, and runtime security controls.

## Implemented Foundation

- Go API platform (`cmd/api`) with auth/RBAC and evidence-chain workflows
- Operator CLI (`cmd/cli`) for evidence create/verify, audit query, enrichment trigger
- OpenAPI spec (`openapi/openapi.yaml`) and MCP action scaffold (`mcp/manifest.json`)
- Unit/API tests (`go test ./...`)
- Robot Framework E2E and OWASP Top 10 suites
- Runtime security policy artifacts (Tetragon profiles + event mapping)
- CI/CD workflows for lint/test/security/release/artifacts/changelog

## API Endpoints

- `GET /healthz`, `GET /readyz`, `GET /version`
- `POST /v1/evidence`, `GET /v1/evidence/:id`
- `POST /v1/custody`, `GET /v1/chain/:id`
- `POST /v1/attest`
- `POST /v1/verify/:id`, `POST /v1/verify/bulk`
- `GET /v1/search`, `GET /v1/audit`
- `POST /v1/enrich`, `GET /v1/enrich/:job_id`
- `GET /v1/proof/:id`
- `GET /v1/admin/key-rotation`

Auth: `X-API-Key` or bearer token roles (`jwt-admin`, `jwt-operator`, `jwt-viewer`).

## Local Run

```bash
make bootstrap
make run
```

## CLI Usage

```bash
go run ./cmd/cli evidence-create EXT-1 sat-a imagery file=a.tif

go run ./cmd/cli evidence-verify <evidence_id>
go run ./cmd/cli audit-query
go run ./cmd/cli enrich-trigger <evidence_id>
```

Optional env:
- `EVIDENCE_API_URL` (default `http://localhost:8080`)
- `EVIDENCE_API_KEY` (default `dev-api-key`)

## Docker Dev

```bash
make dev-up
make dev-down
```

## Testing

```bash
make test
make robot
make owasp
make security-test
```

## Roadmap

Detailed 300-commit implementation roadmap:
- `docs/IMPLEMENTATION_PLAN_300_COMMITS.md`
