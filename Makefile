# Makefile

# Variables
GO := go
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
STATICCHECK := $(shell go env GOPATH)/bin/staticcheck

# Default Go version
GO_VERSION := 1.22

# Targets
.PHONY: all lint golangci-lint staticcheck install-tools

all: lint

lint: golangci-lint staticcheck

golangci-lint:
	$(GOLANGCI_LINT) run --timeout 5m ./...

staticcheck:
	$(STATICCHECK) ./...

install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
