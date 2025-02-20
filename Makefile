OS = $(shell uname | tr A-Z a-z)

PROJ=holos
ORG_PATH=github.com/holos-run
REPO_PATH=$(ORG_PATH)/$(PROJ)

VERSION := $(shell cat version/embedded/major version/embedded/minor version/embedded/patch | xargs printf "%s.%s.%s")
BIN_NAME := holos

DOCKER_REPO=quay.io/holos-run/holos
IMAGE_NAME=$(DOCKER_REPO)

$( shell mkdir -p bin)

# For buf plugin protoc-gen-connect-es
export PATH := $(PWD)/internal/frontend/holos/node_modules/.bin:$(PATH)

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_SUFFIX=$(shell test -n "`git status --porcelain`" && echo "-dirty" || echo "")
GIT_DETAIL=$(shell git describe --tags HEAD)
GIT_TREE_STATE=$(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE=$(shell date -Iseconds)

LD_FLAGS="-w -X ${ORG_PATH}/${PROJ}/version.GitDescribe=${GIT_DETAIL}${GIT_SUFFIX} -X ${ORG_PATH}/${PROJ}/version.GitCommit=${GIT_COMMIT} -X ${ORG_PATH}/${PROJ}/version.GitTreeState=${GIT_TREE_STATE} -X ${ORG_PATH}/${PROJ}/version.BuildDate=${BUILD_DATE}"

.PHONY: default
default: test

.PHONY: bump
bump: bumppatch

.PHONY: bumppatch
bumppatch: ## Bump the patch version.
	scripts/bump patch
	HOLOS_UPDATE_SCRIPTS=1 scripts/test

.PHONY: bumpminor
bumpminor: ## Bump the minor version.
	scripts/bump minor
	scripts/bump patch 0
	HOLOS_UPDATE_SCRIPTS=1 scripts/test

.PHONY: bumpmajor
bumpmajor: ## Bump the major version.
	scripts/bump major
	scripts/bump minor 0
	scripts/bump patch 0
	HOLOS_UPDATE_SCRIPTS=1 scripts/test

.PHONY: show-version
show-version: ## Print the full version.
	@echo $(VERSION)

.PHONY: tag
tag: ## Tag a release
	git tag v$(VERSION)

.PHONY: tidy
tidy: ## Tidy go module.
	go mod tidy

.PHONY: fmt
fmt: ## Format code.
	cd internal/generate/platforms && cue fmt ./...
	go fmt ./...

.PHONY: vet
vet: ## Vet Go code.
	go vet ./...

.PHONY: generate
generate: ## Generate code.
	go generate ./...

.PHONY: build
build: ## Build holos executable.
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -trimpath -o bin/$(BIN_NAME) -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/$(BIN_NAME)

linux: ## Build holos executable for tilt.
	@echo "building ${BIN_NAME}.linux ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	GOOS=linux go build -trimpath -o bin/$(BIN_NAME).linux -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/$(BIN_NAME)

.PHONY: install
install: build ## Install holos to GOPATH/bin
	install bin/$(BIN_NAME) $(shell go env GOPATH)/bin/$(BIN_NAME)

.PHONY: clean
clean: ## Clean executables.
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

.PHONY: test
test: ## Run tests.
	scripts/test

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: lint
lint: golangci-lint ## Run linters.
	buf lint
	cd internal/frontend/holos && ng lint
	./hack/cspell

.PHONY: coverage
coverage: test  ## Test coverage profile.
	go tool cover   -html=coverage.out

.PHONY: snapshot
snapshot:  ## Go release snapshot
	goreleaser release --snapshot --clean

.PHONY: tools
tools: go-deps frontend-deps website-deps ## install tool dependencies

.PHONY: go-deps
go-deps: ## tool versions pinned in tools.go
	go install cuelang.org/go/cmd/cue
	go install github.com/bufbuild/buf/cmd/buf
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go
	go install honnef.co/go/tools/cmd/staticcheck
	go install golang.org/x/tools/cmd/godoc
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc
	go install github.com/google/ko
	# curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash

.PHONY: frontend-deps
frontend-deps: ## Install Angular deps for go generate
	cd internal/frontend/holos && npm install

.PHONY: website-deps
website-deps: ## Install Docusaurus deps for go generate
	cd doc/website && npm install

.PHONY: image # refer to .ko.yaml as well
image:  ## Container image build for workflows/publish.yaml
	KO_DOCKER_REPO=$(DOCKER_REPO) GIT_DETAIL=$(GIT_DETAIL) GIT_SUFFIX=$(GIT_SUFFIX) ko build --platform=all --bare ./cmd/holos --tags $(GIT_DETAIL)$(GIT_SUFFIX) --tags latest

.PHONY: prod-deploy
prod-deploy: install image  ## deploy to PROD
	GIT_DETAIL=$(GIT_DETAIL) GIT_SUFFIX=$(GIT_SUFFIX) bash ./hack/deploy

.PHONY: dev-deploy
dev-deploy: install image  ## deploy to dev
	GIT_DETAIL=$(GIT_DETAIL) GIT_SUFFIX=$(GIT_SUFFIX) bash ./hack/deploy-dev

.PHONY: website
website: ## Build website
	./hack/build-website

.PHONY: unity
unity: ## https://cuelabs.dev/unity/
	./scripts/unity

.PHONY: update-docs
update-docs: ## Update doc examples
	HOLOS_UPDATE_SCRIPTS=1 go test -v ./doc/md/...

.PHONY: help
help:  ## Display this help menu.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
