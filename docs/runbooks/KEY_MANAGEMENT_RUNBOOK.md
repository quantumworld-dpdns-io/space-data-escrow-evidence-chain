# Key Management Runbook

## Scope
Covers API key rotation metadata and attestation signature key lifecycle.

## Rotation Steps
1. Generate new key pair in secure KMS/HSM workflow.
2. Update active key metadata endpoint backing values.
3. Validate `GET /v1/admin/key-rotation` returns expected next rotation window.
4. Run dual-sign mode checks for transition period.

## Emergency Rotation
1. Revoke compromised key immediately.
2. Issue replacement and notify operators.
3. Re-sign critical attestations where required.
