.ONESHELL:
.DELETE_ON_ERROR:
export SHELL     := bash
export SHELLOPTS := pipefail:errexit
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rule

# Adapted from https://suva.sh/posts/well-documented-makefiles/
.PHONY: help
help: ## Display this help
help:
	@awk 'BEGIN {FS = ": ##"; printf "Usage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_\.\-\/%]+: ##/ { printf "  %-45s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.linted: cmd/fit
	if golangci-lint run $<; then touch $@; fi

lint: ## Lint Go source code.
lint: .linted

vendor: ## Update vendored Go source code.
vendor: go.mod go.sum
	go mod tidy && go mod vendor

result/bin/fit: ## Build the fit binary using Nix.
result/bin/fit: vendor default.nix flake.nix lint
	nix build .

pgsql: ## Generate database code.
pgsql: sqlc.json query.sql schema.sql
	sqlc generate
