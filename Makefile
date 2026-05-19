.PHONY: run test robot owasp test-all security-test dev-up dev-down build cli bootstrap

run:
	APP_API_KEY=dev-api-key go run ./cmd/api

build:
	go build ./cmd/api

cli:
	go build ./cmd/cli

test:
	go test ./...

robot:
	robot tests/robot/suites

owasp:
	robot tests/robot/owasp

test-all: test robot owasp

security-test:
	API_KEY=dev-api-key bash scripts/security/dast_lite.sh http://localhost:8080
	robot tests/robot/owasp

bootstrap:
	bash scripts/dev/bootstrap_local.sh

dev-up:
	docker compose up -d --build

dev-down:
	docker compose down --remove-orphans
