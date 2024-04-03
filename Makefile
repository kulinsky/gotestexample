.PHONY: lint
default: run

COVERAGE_FILE ?= coverage.out

.PHONY: lint
lint:
	@golangci-lint --version && echo "golangci-lint -v run --fix ./..." || echo "golangci-lint not found"
	@golangci-lint -v run --fix ./...

.PHONY: run
run:
	go run ./cmd/server

.PHONY: test
test:
	@go test -v -coverpkg='github.com/kulinsky/gotestexample/...' -race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'

.PHONY: test-short
test-short:
	@go test -short -v ./...

.PHONY: mockgen
mockgen:
	@go generate ./...
