
.PHONY: default
default:  build

BIN_DIR := $(shell pwd)/bin

# LDFLAGS
VERSION := $(shell git tag --sort=committerdate | tail -1)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u '+%Y-%m-%d')
LDFLAGS=-ldflags \
				" \
				-X github.com/mirantiscontainers/boundless-cli/cmd.version=${VERSION} \
				-X github.com/mirantiscontainers/boundless-cli/cmd.commit=${COMMIT} \
				-X github.com/mirantiscontainers/boundless-cli/cmd.date=${DATE} \
				"

.PHONY: build
build:  ## build locally
	@go mod download
	@CGO_ENABLED=0 go build ${LDFLAGS} -o ${BIN_DIR}/bctl ./

.PHONY: install
install:  ## install locally
	@go mod download
	@CGO_ENABLED=0 go build ${LDFLAGS} -o ${GOPATH}/bin/bctl ./

.PHONY: init
init:
	@${BIN_DIR}/bctl init

.PHONY: apply
apply:
	@${BIN_DIR}/bctl apply --config blueprint.yaml

.PHONY: update
update:
	@${BIN_DIR}/bctl update --config blueprint.yaml

.PHONY: reset
reset:
	@${BIN_DIR}/bctl reset --config blueprint.yaml

.PHONY: build-charts
build-charts:
	@cd ./charts && make build

.PHONY: test
test:  ## Run tests.
	@go test ./... -coverprofile cover.out

.PHONY: vet
vet: ## Run go vet against code.
	@go vet ./...
