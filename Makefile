
.PHONY: default
default:  build

BIN_DIR := $(shell pwd)/bin
VERSION := dev-$(shell git rev-parse --short HEAD)

.PHONY: help
help: ## Display this help output
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: build
build:  ## Build the binary
	@CGO_ENABLED=0 go build -ldflags "-X 'boundless-cli/cmd.version=${VERSION}'" -o ${BIN_DIR}/bctl ./


.PHONY: clean
clean:  ## Clean up build artifacts
	@rm -rf bin/
	@rm -rf cover.out

.PHONY: install
install:  ## Install the binary on your system
	@CGO_ENABLED=0 go build -ldflags "-X 'boundless-cli/cmd.version=${VERSION}'" -o ${GOPATH}/bin/bctl ./

.PHONY: build-charts
build-charts: ## Build the charts
	@cd ./charts && make build

##@ Testing

.PHONY: test
test:  ## Run unit tests
# Skip the test directory since it containers integration and e2e tests
	@go test $$(go list ./... | grep -v /test/) -coverprofile cover.out

.PHONY: vet
vet: ## Run go vet against code.
	@go vet ./...


##@ Binary Commands

.PHONY: init
init: ## Initialize a default blueprint
	@${BIN_DIR}/bctl init > blueprint.yaml

.PHONY: apply
apply: ## Apply the blueprint
	@${BIN_DIR}/bctl apply --config blueprint.yaml

.PHONY: update
update: ## Update the cluster with the blueprint
	@${BIN_DIR}/bctl update --config blueprint.yaml

.PHONY: reset
reset: ## Reset the cluster (remove all resources)
	@${BIN_DIR}/bctl reset --config blueprint.yaml


