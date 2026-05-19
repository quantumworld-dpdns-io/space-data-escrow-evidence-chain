#!/usr/bin/env bash
set -euo pipefail

required=(
  "docs/runbooks/RELEASE_CANDIDATE_CHECKLIST.md"
  "docs/runbooks/GO_LIVE_CHECKLIST.md"
  "docs/runbooks/GOVERNANCE_SIGNOFF.md"
  "docs/runbooks/V1_RELEASE_PROTOCOL.md"
  "docs/runbooks/V1_ACCEPTANCE_MATRIX.md"
  "docs/runbooks/V1_SIGNOFF_BUNDLE_TEMPLATE.md"
)

for f in "${required[@]}"; do
  [[ -f "$f" ]] || { echo "missing required doc: $f"; exit 1; }
done

echo "v1 docs validation passed"
