#!/usr/bin/env bash
set -euo pipefail
BASE_URL="${1:-http://localhost:8080}"
API_KEY="${API_KEY:-dev-api-key}"

check() {
  local method="$1"
  local path="$2"
  local expected="$3"
  local code
  code=$(curl -s -o /tmp/dast_body.txt -w "%{http_code}" -X "$method" -H "X-API-Key: ${API_KEY}" "${BASE_URL}${path}")
  if [[ "$code" != "$expected" ]]; then
    echo "DAST check failed: ${method} ${path} expected ${expected} got ${code}"
    cat /tmp/dast_body.txt || true
    exit 1
  fi
}

check GET "/healthz" 200
check GET "/readyz" 200
check GET "/v1/search?page=1&page_size=1" 200
check POST "/v1/verify/does-not-exist" 404

echo "DAST-lite checks passed"
