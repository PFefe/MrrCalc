# Makefile

# Variables
GO := go
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
GOROOT := /usr/local/go

# Default Go version
GO_VERSION := 1.22

# Targets
.PHONY: all lint golangci-lint install-tools

all: lint

lint: golangci-lint

golangci-lint:
    @GOROOT=$(GOROOT) $(GOLANGCI_LINT) run --timeout 5m ./...

install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1
