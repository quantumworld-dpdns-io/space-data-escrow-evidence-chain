# Operations Runbook

## Service Startup
- Local: `make run`
- Container: `make dev-up`

## Health Checks
- Liveness: `GET /healthz`
- Readiness: `GET /readyz`
- Version: `GET /version`

## Routine Checks
1. Validate API auth (`/v1/search` should return 401 without credentials).
2. Run `make test` and `make security-test` before deploy.
3. Verify audit stream (`GET /v1/audit`) includes recent events.

## Deployment Verification
1. Create evidence via API/CLI.
2. Append custody + verify chain.
3. Confirm no regression in latency benchmark (`make benchmark`).

## Rollback
1. Roll back to previous image/tag.
2. Re-run health checks.
3. Re-run a minimal evidence lifecycle scenario.
