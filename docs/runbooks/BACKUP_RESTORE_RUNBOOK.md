# Backup and Restore Runbook

## Backup Scope
- Evidence records
- Custody events
- Attestations
- Audit entries

## Strategy
1. Export proof bundles for critical evidence IDs.
2. Snapshot metadata store (DuckDB target state).
3. Archive runtime security events.

## Restore Verification
1. Restore datasets into staging.
2. Run chain verification on representative samples.
3. Compare checksum/hash integrity and attestation counts.
