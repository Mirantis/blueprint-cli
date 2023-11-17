
.PHONY: default
default:  build

BIN_DIR := $(shell pwd)/bin
VERSION := dev-$(shell git rev-parse --short HEAD)

.PHONY: build
build:  ## build locally
	@go mod download
	@CGO_ENABLED=0 go build -ldflags "-X 'boundless-cli/cmd.version=${VERSION}'" -o ${BIN_DIR}/bctl ./

.PHONY: install
install:  ## install locally
	@go mod download
	@CGO_ENABLED=0 go build -ldflags "-X 'boundless-cli/cmd.version=${VERSION}'" -o ${GOPATH}/bin/bctl ./

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
