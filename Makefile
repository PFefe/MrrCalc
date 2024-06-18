# Makefile

# Variables
GO := go
GOLANGCI_LINT := golangci-lint
STATICCHECK := staticcheck

# Default Go version
GO_VERSION := 1.22

# Targets
.PHONY: all lint golangci-lint staticcheck

all: lint

lint: golangci-lint staticcheck

golangci-lint:
	$(GOLANGCI_LINT) run --timeout 5m

staticcheck:
	$(STATICCHECK) ./...

# Install dependencies
install:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
