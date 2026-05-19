# space-data-escrow-evidence-chain

Space data escrow and evidence chain platform with verifiable chain-of-custody APIs, audit trails, and security testing.

## Implemented Foundation

- Go + Gin API platform (`cmd/api`)
- Evidence lifecycle APIs: create/get/search/verify
- Custody event capture and audit trail
- Deterministic payload hashing (SHA-256)
- OpenAPI spec (`openapi/openapi.yaml`)
- Unit/API tests (`go test ./...`)
- Robot Framework E2E and OWASP Top 10 baseline suites
- GitHub Actions CI + security workflows

## API Endpoints

- `GET /healthz`
- `GET /readyz`
- `POST /v1/evidence`
- `GET /v1/evidence/:id`
- `POST /v1/custody`
- `POST /v1/verify/:id`
- `GET /v1/search?q=...`
- `GET /v1/audit`

All `/v1/*` endpoints require `X-API-Key`.

## Local Run

```bash
go mod tidy
make run
```

## Test

```bash
make test
make robot
make owasp
```

## Roadmap

Detailed 300-commit implementation roadmap:
- `docs/IMPLEMENTATION_PLAN_300_COMMITS.md`
