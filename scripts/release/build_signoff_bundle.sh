#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="${1:-artifacts/signoff}"
mkdir -p "$OUT_DIR/docs/runbooks" "$OUT_DIR/workflows" "$OUT_DIR/security"

# Ensure required docs are present and current
bash scripts/release/validate_v1_docs.sh
bash scripts/release/validate_v1_dossier.sh

copy_file() {
  local src="$1"
  local dst="$2"
  mkdir -p "$(dirname "$dst")"
  cp "$src" "$dst"
}

# Core runbooks and dossier
for f in \
  docs/runbooks/RELEASE_CANDIDATE_CHECKLIST.md \
  docs/runbooks/GO_LIVE_CHECKLIST.md \
  docs/runbooks/POST_DEPLOY_VERIFICATION.md \
  docs/runbooks/GOVERNANCE_SIGNOFF.md \
  docs/runbooks/V1_ACCEPTANCE_MATRIX.md \
  docs/runbooks/V1_RELEASE_PROTOCOL.md \
  docs/runbooks/V1_SIGNOFF_BUNDLE_TEMPLATE.md \
  docs/runbooks/V1_LAUNCH_DOSSIER.md \
  docs/runbooks/RELEASE_EVIDENCE_INDEX.md
 do
  copy_file "$f" "$OUT_DIR/$f"
 done

# Workflow references
for wf in .github/workflows/ci.yml .github/workflows/security.yml .github/workflows/release-readiness.yml .github/workflows/release.yml .github/workflows/post-deploy-verify.yml; do
  copy_file "$wf" "$OUT_DIR/workflows/$(basename "$wf")"
done

# Runtime security policy + schema
copy_file security/events/runtime_event_schema.json "$OUT_DIR/security/runtime_event_schema.json"
copy_file security/tetragon/policy-suspicious-exec.yaml "$OUT_DIR/security/policy-suspicious-exec.yaml"
copy_file security/tetragon/policy-egress-anomaly.yaml "$OUT_DIR/security/policy-egress-anomaly.yaml"

# Index file for reviewers
cat > "$OUT_DIR/SIGNOFF_BUNDLE_INDEX.md" <<DOC
# Sign-Off Bundle Index

Generated at: $(date -u +%Y-%m-%dT%H:%M:%SZ)

Included:
- Runbooks/checklists and v1 dossier
- Core CI/Security/Release workflows
- Runtime security schema and tetragon policies

Use this bundle alongside CI run links and deployment verification URLs.
DOC

echo "sign-off bundle assembled at $OUT_DIR"
