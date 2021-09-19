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

cmd/db/main.go: $(wildcard pgsql/*.go)
	touch $@

%/.linted: ## Lint source code.
%/.linted: %/main.go
	if golangci-lint run ./$<; then touch $@; fi

vendor: ## Update vendored Go source code.
vendor: go.mod go.sum
	go mod tidy && go mod vendor

.built: ## Build binaries using Nix.
.built: vendor default.nix flake.nix cmd/db/.linted cmd/files/.linted cmd/fit/.linted cmd/settings/.linted
	if nix build .; then touch $@; fi

pgsql/db.go pgsql/models.go pgsql/query.sql.go: ## Generate database code.
pgsql/db.go pgsql/models.go pgsql/query.sql.go: sqlc.json query.sql schema.sql
	sqlc generate

docker-compose.yml: ## Generate the docker-compose from Jsonnet.
docker-compose.yml: docker-compose.jsonnet
	jsonnet $< > $@

VIVOSPORT_DEV ?= /dev/disk/by-label/GARMIN
.PHONY: data
data: ## Rsync data from the vivosport device.
data:
	sudo mount $(VIVOSPORT_DEV) /tmp/garmin
	rsync -avz /tmp/garmin/GARMIN/* $@
