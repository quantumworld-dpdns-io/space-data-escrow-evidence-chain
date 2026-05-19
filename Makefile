.PHONY: run test robot owasp test-all

run:
	APP_API_KEY=dev-api-key go run ./cmd/api

test:
	go test ./...

robot:
	robot tests/robot/suites

owasp:
	robot tests/robot/owasp

test-all: test robot owasp
