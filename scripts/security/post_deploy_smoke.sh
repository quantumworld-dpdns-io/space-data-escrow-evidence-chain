#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:?usage: post_deploy_smoke.sh <base_url>}"
API_KEY="${API_KEY:?API_KEY env is required}"

check() {
  local method="$1"
  local path="$2"
  local expected="$3"
  local code
  code=$(curl -s -o /tmp/postdeploy_body.txt -w "%{http_code}" -X "$method" -H "X-API-Key: ${API_KEY}" "${BASE_URL}${path}")
  if [[ "$code" != "$expected" ]]; then
    echo "post-deploy check failed: ${method} ${path} expected ${expected} got ${code}"
    cat /tmp/postdeploy_body.txt || true
    exit 1
  fi
}

# unauthenticated checks
code=$(curl -s -o /tmp/postdeploy_body.txt -w "%{http_code}" "${BASE_URL}/healthz")
[[ "$code" == "200" ]] || { echo "healthz failed with ${code}"; exit 1; }
code=$(curl -s -o /tmp/postdeploy_body.txt -w "%{http_code}" "${BASE_URL}/readyz")
[[ "$code" == "200" ]] || { echo "readyz failed with ${code}"; exit 1; }

# authenticated checks
check GET "/v1/search?page=1&page_size=1" 200
check POST "/v1/verify/does-not-exist" 404
check GET "/v1/audit" 200

echo "post-deploy smoke checks passed"
