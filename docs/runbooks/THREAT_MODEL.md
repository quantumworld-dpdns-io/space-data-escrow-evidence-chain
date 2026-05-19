# Threat Model

## Assets
- Evidence payload integrity
- Chain-of-custody ordering
- Attestation signatures
- Audit event trail
- Runtime security event logs

## Primary Threats
- Unauthorized data mutation
- Replay/forgery of attestations
- Privilege abuse on admin endpoints
- Supply-chain/runtime compromise

## Mitigations
- Auth + RBAC with explicit role boundaries
- Hash-based integrity + custody monotonic checks
- OWASP A01-A10 test coverage + DAST-lite + ZAP baseline
- Runtime event ingestion and policy-based detection
