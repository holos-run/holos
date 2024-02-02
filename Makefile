OS = $(shell uname | tr A-Z a-z)

PROJ=holos
ORG_PATH=github.com/holos-run
REPO_PATH=$(ORG_PATH)/$(PROJ)

VERSION := $(shell grep "const Version " pkg/version/version.go | sed -E 's/.*"(.+)"$$/\1/')
BIN_NAME := holos

DOCKER_REPO=quay.io/openinfrastructure/holos
IMAGE_NAME=$(DOCKER_REPO)

$( shell mkdir -p bin)

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_TREE_STATE=$(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE=$(shell date -Iseconds)

LD_FLAGS="-w -X ${ORG_PATH}/${PROJ}/pkg/version.GitCommit=${GIT_COMMIT} -X ${ORG_PATH}/${PROJ}/pkg/version.GitTreeState=${GIT_TREE_STATE} -X ${ORG_PATH}/${PROJ}/pkg/version.BuildDate=${BUILD_DATE}"

.PHONY: default
default: test

.PHONY: bump
bump: bumppatch

.PHONY: bumppatch
bumppatch: ## Bump the patch version.
	scripts/bump patch

.PHONY: bumpminor
bumpminor: ## Bump the minor version.
	scripts/bump minor
	scripts/bump patch 0

.PHONY: bumpmajor
bumpmajor: ## Bump the major version.
	scripts/bump major
	scripts/bump minor 0
	scripts/bump patch 0

.PHONY: tidy
tidy: ## Tidy go module.
	go mod tidy

.PHONY: fmt
fmt: ## Format Go code.
	go fmt ./...

.PHONY: vet
vet: ## Vet Go code.
	go vet ./...

.PHONY: generate
generate: ## Generate code.
	go generate ./...

.PHONY: build
build: generate ## Build holos executable.
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -trimpath -o bin/$(BIN_NAME) -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/$(BIN_NAME)

.PHONY: install
install: build ## Install holos to GOPATH/bin
	install bin/$(BIN_NAME) $(shell go env GOPATH)/bin/$(BIN_NAME)

.PHONY: clean
clean: ## Clean executables.
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

.PHONY: test
test: ## Run tests.
	go test -coverpkg=./... -coverprofile=coverage.out ./...

.PHONY: coverage
coverage: test  ## Test coverage profile.
	go tool cover   -html=coverage.out

.PHONY: snapshot
snapshot:  ## Go release snapshot
	goreleaser release --snapshot --clean

.PHONY: help
help:  ## Display this help menu.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
