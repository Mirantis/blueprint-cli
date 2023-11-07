
.PHONY: default
default:  build

BIN_DIR := $(shell pwd)/bin

.PHONY: build
build:  ## build locally
	@go mod download
	@CGO_ENABLED=1 go build -o ${BIN_DIR}/bctl ./

.PHONY: install
install:  ## install locally
	@go mod download
	@CGO_ENABLED=1 go build -o ${GOPATH}/bin/bctl ./

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
