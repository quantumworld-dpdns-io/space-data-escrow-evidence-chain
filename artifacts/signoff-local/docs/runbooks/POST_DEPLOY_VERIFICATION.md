# Post-Deploy Verification

## Automated
Run workflow: `Post-Deploy Verify` with target URL and API key.

Automated checks include:
- `GET /healthz`
- `GET /readyz`
- authenticated search
- verify-not-found path behavior

## Manual
1. Create evidence with CLI.
2. Append custody and verify chain.
3. Trigger enrichment and fetch job status.
4. Confirm audit entries include deploy-time actions.
