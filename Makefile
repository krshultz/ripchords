.PHONY: build test vet lint check tools hooks

# Pinned to match CI (.github/workflows/ci.yml) so local and CI never disagree.
GOLANGCI_VERSION := v2.12.2

build:
	go build -ldflags "-X main.version=$(shell git describe --tags --always)" -o ripchords .

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

# Run everything CI runs.
check: vet lint test

# Install developer tooling (golangci-lint) into $(go env GOPATH)/bin.
tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)

# Enable the repo's git hooks (runs checks before each commit).
hooks:
	git config core.hooksPath .githooks
	@echo "git hooks enabled (core.hooksPath = .githooks)"
