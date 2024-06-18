# Makefile

# Variables
GO := go
GOLANGCI_LINT := golangci-lint
STATICCHECK := staticcheck

# Default Go version
GO_VERSION := 1.22
GOROOT_PATH := /usr/local/go
GOPATH_PATH := /Users/efe/go

# Targets
.PHONY: all lint golangci-lint staticcheck

all: lint

lint: export GOROOT=$(GOROOT_PATH)
lint: export GOPATH=$(GOPATH_PATH)
lint: golangci-lint staticcheck

golangci-lint: export GOROOT=$(GOROOT_PATH)
golangci-lint: export GOPATH=$(GOPATH_PATH)
golangci-lint:
	$(GOLANGCI_LINT) run --timeout 5m

staticcheck: export GOROOT=$(GOROOT_PATH)
staticcheck: export GOPATH=$(GOPATH_PATH)
staticcheck:
	$(STATICCHECK) ./...

# Install dependencies
install:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
