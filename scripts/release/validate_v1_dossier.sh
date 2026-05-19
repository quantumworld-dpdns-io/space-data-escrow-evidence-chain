#!/usr/bin/env bash
set -euo pipefail

TARGET="docs/runbooks/V1_LAUNCH_DOSSIER.md"
TMP="/tmp/V1_LAUNCH_DOSSIER.generated.md"

# Generate fresh content to temp by reusing generator logic but overriding output path
DATE_UTC="$(date -u +%Y-%m-%d)"
cat > "$TMP" <<DOC
# v1 Launch Dossier

Generated on: ${DATE_UTC} (UTC)

## 1. Release Identity
- Target release: \`v1.0.0\`
- Repository: \`space-data-escrow-evidence-chain\`
- Protocol: [V1 Release Protocol](./V1_RELEASE_PROTOCOL.md)

## 2. Mandatory Checklists
- [Release Candidate Checklist](./RELEASE_CANDIDATE_CHECKLIST.md)
- [Go-Live Checklist](./GO_LIVE_CHECKLIST.md)
- [Post-Deploy Verification](./POST_DEPLOY_VERIFICATION.md)
- [Governance Sign-Off](./GOVERNANCE_SIGNOFF.md)
- [v1 Acceptance Matrix](./V1_ACCEPTANCE_MATRIX.md)
- [v1 Sign-Off Bundle Template](./V1_SIGNOFF_BUNDLE_TEMPLATE.md)

## 3. Required Workflow Evidence
- [CI Workflow](../../.github/workflows/ci.yml)
- [Security Workflow](../../.github/workflows/security.yml)
- [Release Readiness Workflow](../../.github/workflows/release-readiness.yml)
- [Release Workflow](../../.github/workflows/release.yml)
- [Post-Deploy Verify Workflow](../../.github/workflows/post-deploy-verify.yml)
- [ZAP Baseline Workflow](../../.github/workflows/zap-baseline.yml)

## 4. Artifact Evidence
- Binary artifact: \`artifacts/evidence-api\`
- SBOM: \`artifacts/sbom.spdx.json\`
- Provenance: \`artifacts/provenance.json\`
- Release notes: \`artifacts/release-notes.md\`
- Changelog draft artifact from workflow

## 5. Security and Runtime Controls
- Runtime event schema: [runtime_event_schema.json](../../security/events/runtime_event_schema.json)
- Tetragon suspicious exec policy: [policy-suspicious-exec.yaml](../../security/tetragon/policy-suspicious-exec.yaml)
- Tetragon egress anomaly policy: [policy-egress-anomaly.yaml](../../security/tetragon/policy-egress-anomaly.yaml)
- OWASP suite: [owasp_top10.robot](../../tests/robot/owasp/owasp_top10.robot)
- DAST-lite script: [dast_lite.sh](../../scripts/security/dast_lite.sh)

## 6. Operational Readiness
- [Operations Runbook](./OPERATIONS_RUNBOOK.md)
- [Incident Response Runbook](./INCIDENT_RESPONSE_RUNBOOK.md)
- [Key Management Runbook](./KEY_MANAGEMENT_RUNBOOK.md)
- [Backup and Restore Runbook](./BACKUP_RESTORE_RUNBOOK.md)
- [Retention and Compliance Policy](./RETENTION_COMPLIANCE_POLICY.md)
- [Threat Model](./THREAT_MODEL.md)
- [SLO and Alerts](./SLO_AND_ALERTS.md)

## 7. Approval Record (to complete during release)
- Release tag:
- Commit SHA:
- CI run links:
- Security run links:
- Post-deploy verification link:
- Engineering approval:
- Security approval:
- Operations approval:
DOC

if ! diff -u "$TMP" "$TARGET" >/tmp/v1_dossier_diff.txt; then
  echo "v1 dossier drift detected; run scripts/release/generate_v1_dossier.sh"
  cat /tmp/v1_dossier_diff.txt
  exit 1
fi

echo "v1 dossier validation passed"
