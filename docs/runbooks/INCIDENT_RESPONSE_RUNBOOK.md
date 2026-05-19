# Incident Response Runbook

## Trigger Conditions
- Elevated 5xx responses
- Security alerts from runtime event ingestion
- DAST/OWASP failures in CI

## Immediate Actions
1. Freeze production deploys.
2. Capture logs and `/v1/audit` output.
3. Identify blast radius (tenant/source/type scope).

## Containment
1. Rotate API keys if compromise suspected.
2. Restrict high-risk endpoints (`/v1/admin/*`, `/v1/attest`) via gateway/WAF.
3. Enable enhanced runtime monitoring profile.

## Recovery
1. Patch root cause.
2. Validate with unit + Robot + OWASP + DAST-lite.
3. Deploy with staged verification.
