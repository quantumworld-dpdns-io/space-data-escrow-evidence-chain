#!/usr/bin/env bash
set -euo pipefail

echo "[bootstrap] go mod tidy"
go mod tidy

echo "[bootstrap] local tools"
python3 -m pip install --user robotframework robotframework-requests >/dev/null 2>&1 || true

echo "[bootstrap] done"
