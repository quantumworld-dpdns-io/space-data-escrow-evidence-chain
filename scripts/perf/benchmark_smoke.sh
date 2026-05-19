#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
API_KEY="${API_KEY:-dev-api-key}"
ITERATIONS="${ITERATIONS:-20}"
MAX_P95_MS="${MAX_P95_MS:-500}"

measure() {
  local url="$1"
  local i=0
  local times=()
  while [ "$i" -lt "$ITERATIONS" ]; do
    local t
    t=$(curl -s -o /dev/null -w "%{time_total}" -H "X-API-Key: ${API_KEY}" "${url}")
    # convert seconds float to ms integer
    local ms
    ms=$(awk -v v="$t" 'BEGIN{printf "%d", v*1000}')
    times+=("$ms")
    i=$((i+1))
  done
  printf "%s\n" "${times[@]}" | sort -n > /tmp/bench_times.txt
  local idx=$(( (95 * ITERATIONS + 99) / 100 - 1 ))
  if [ "$idx" -lt 0 ]; then idx=0; fi
  local p95
  p95=$(sed -n "$((idx+1))p" /tmp/bench_times.txt)
  echo "$p95"
}

P95_SEARCH=$(measure "${BASE_URL}/v1/search?page=1&page_size=5")
P95_HEALTH=$(measure "${BASE_URL}/healthz")

echo "benchmark.p95.search_ms=${P95_SEARCH}"
echo "benchmark.p95.health_ms=${P95_HEALTH}"

if [ "$P95_SEARCH" -gt "$MAX_P95_MS" ]; then
  echo "p95 search latency regression: ${P95_SEARCH}ms > ${MAX_P95_MS}ms"
  exit 1
fi

echo "benchmark smoke passed"
