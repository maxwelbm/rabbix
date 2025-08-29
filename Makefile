SHELL := /usr/bin/env bash

# Versão do golangci-lint compatível e fixada para reprodutibilidade
GOLANGCI_LINT_VERSION ?= v2.1.6

.PHONY: all lint golangci-lint-install

all:
	go build -o rabbix main.go

golangci-lint-install:
	@set -euo pipefail; \
	if ! command -v golangci-lint >/dev/null 2>&1 || ! golangci-lint version 2>/dev/null | grep -q "version $(GOLANGCI_LINT_VERSION)"; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	else \
		echo "golangci-lint $(GOLANGCI_LINT_VERSION) already installed"; \
	fi

lint: golangci-lint-install
	@set -euo pipefail; \
	echo "Running linters with golangci-lint $(GOLANGCI_LINT_VERSION)"; \
	golangci-lint run ./...