#!/usr/bin/env bash
set -euo pipefail
mkdir -p artifacts
cat > artifacts/sbom.spdx.json <<JSON
{
  "spdxVersion": "SPDX-2.3",
  "name": "space-data-escrow-evidence-chain",
  "creationInfo": {
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "creators": ["Tool: local-sbom-placeholder"]
  }
}
JSON
