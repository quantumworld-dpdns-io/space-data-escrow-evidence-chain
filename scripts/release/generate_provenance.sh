#!/usr/bin/env bash
set -euo pipefail
mkdir -p artifacts
cat > artifacts/provenance.json <<JSON
{
  "project": "space-data-escrow-evidence-chain",
  "generated_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "builder": "github-actions-placeholder"
}
JSON
