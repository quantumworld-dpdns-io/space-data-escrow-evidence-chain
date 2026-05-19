# v1 Acceptance Matrix

| Area | Criteria | Validation Method | Status |
|---|---|---|---|
| API Core | Evidence create/get/search/verify works | Unit + Robot | Pending |
| Chain Integrity | Hash + custody monotonic checks active | Unit tests | Pending |
| Security AuthZ | RBAC denies unauthorized admin access | API tests + OWASP A01 | Pending |
| Runtime Security | Runtime event ingestion + Tetragon policy files present | Service tests + file checks | Pending |
| Compliance | Audit endpoint and retention policy docs present | Robot + docs review | Pending |
| Performance | Benchmark regression gate passes | CI benchmark job | Pending |
| Release | Artifacts (binary/SBOM/provenance) generated | Release/Artifacts workflows | Pending |
| Post-Deploy | Smoke verification workflow passes | Post-Deploy Verify workflow | Pending |
