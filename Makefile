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
%/.linted: %/main.go vendor
	if golangci-lint run ./$<; then touch $@; fi

vendor: ## Update vendored Go source code.
vendor: 21.40.xlsx go.mod go.sum
	go mod tidy && go mod vendor
	fitgen -verbose -sdk 21.40 $< $@/github.com/tormoder/fit

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


SDK_VERSION ?= 21.60
sdk.zip:
	curl -Lo sdk.zip https://developer.garmin.com/downloads/fit/sdk/FitSDKRelease_$(SDK_VERSION).00.zip

sdk: ## Fetch and extract the Garmin SDK.
sdk: sdk.zip
	unzip $< -d $@

csv/%.csv: ## Convert a FIT file to CSV.
csv/%.csv: data/%.FIT sdk
	mkdir -p $(@D)
	java -jar sdk/java/FitCSVTool.jar -b $< $@
