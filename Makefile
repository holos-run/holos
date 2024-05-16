OS = $(shell uname | tr A-Z a-z)

PROJ=holos
ORG_PATH=github.com/holos-run
REPO_PATH=$(ORG_PATH)/$(PROJ)

VERSION := $(shell cat version/embedded/major version/embedded/minor version/embedded/patch | xargs printf "%s.%s.%s")
BIN_NAME := holos

DOCKER_REPO=quay.io/openinfrastructure/holos
IMAGE_NAME=$(DOCKER_REPO)

$( shell mkdir -p bin)

# For buf plugin protoc-gen-connect-es
export PATH := $(PWD)/internal/frontend/holos/node_modules/.bin:$(PATH)

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_TREE_STATE=$(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE=$(shell date -Iseconds)

LD_FLAGS="-w -X ${ORG_PATH}/${PROJ}/version.GitCommit=${GIT_COMMIT} -X ${ORG_PATH}/${PROJ}/version.GitTreeState=${GIT_TREE_STATE} -X ${ORG_PATH}/${PROJ}/version.BuildDate=${BUILD_DATE}"

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

.PHONY: show-version
show-version: ## Print the full version.
	@echo $(VERSION)

.PHONY: tidy
tidy: ## Tidy go module.
	go mod tidy

.PHONY: fmt
fmt: ## Format code.
	cd docs/examples && cue fmt ./...
	go fmt ./...

.PHONY: vet
vet: ## Vet Go code.
	go vet ./...

.PHONY: gencue
gencue: ## Generate CUE definitions
	cd docs/examples && cue get go github.com/holos-run/holos/api/...

.PHONY: rmgen
rmgen: ## Remove generated code
	git rm -rf service/gen/ internal/frontend/holos/src/app/gen/ || true
	rm -rf service/gen/ internal/frontend/holos/src/app/gen/
	git rm -rf internal/ent/
	rm -rf internal/ent/
	git restore --staged internal/ent/generate.go internal/ent/schema/
	git restore internal/ent/generate.go internal/ent/schema/

.PHONY: regenerate
regenerate: generate ## Re-generate code (delete and re-create)

.PHONY: generate
generate: buf ## Generate code.
	go generate ./...

.PHONY: build
build: generate frontend ## Build holos executable.
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
	scripts/test

.PHONY: lint
lint: ## Run linters.
	buf lint
	cd internal/frontend/holos && ng lint
	golangci-lint run

.PHONY: coverage
coverage: test  ## Test coverage profile.
	go tool cover   -html=coverage.out

.PHONY: snapshot
snapshot:  ## Go release snapshot
	goreleaser release --snapshot --clean

.PHONY: buf
buf: ## buf generate
	cd service && buf dep update
	buf generate

.PHONY: tools
tools: go-deps frontend-deps  ## install tool dependencies

.PHONY: go-deps
go-deps: ## tool versions pinned in tools.go
	go install github.com/bufbuild/buf/cmd/buf
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go
	go install honnef.co/go/tools/cmd/staticcheck@latest
	# curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash

.PHONY: frontend-deps
frontend-deps: ## Setup npm and vite
	cd internal/frontend/holos && npm install
	cd internal/frontend/holos && npm install --save-dev @bufbuild/buf @connectrpc/protoc-gen-connect-es
	cd internal/frontend/holos && npm install @connectrpc/connect @connectrpc/connect-web @bufbuild/protobuf
	# https://github.com/connectrpc/connect-query-es/blob/1350b6f07b6aead81793917954bdb1cc3ce09df9/packages/protoc-gen-connect-query/README.md?plain=1#L23
	cd internal/frontend/holos && npm install --save-dev @connectrpc/protoc-gen-connect-query @bufbuild/protoc-gen-es
	cd internal/frontend/holos && npm install @connectrpc/connect-query @bufbuild/protobuf


.PHONY: frontend
frontend: buf
	cd internal/frontend/holos && rm -rf dist
	mkdir -p internal/frontend/holos/dist
	cd internal/frontend/holos && ng build
	touch internal/frontend/frontend.go

.PHONY: help
help:  ## Display this help menu.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
