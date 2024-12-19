OS = $(shell uname | tr A-Z a-z)

PROJ=holos
ORG_PATH=github.com/holos-run
REPO_PATH=$(ORG_PATH)/$(PROJ)

.PHONY: default
default: test

.PHONY: install
install: ## Install holos to GOPATH/bin
	go install github.com/holos-run/holos/cmd/holos

.PHONY: test
test: ## Run go test
	go test

.PHONY: unity
unity: ## https://cuelabs.dev/unity/
	./scripts/unity

.PHONY: help
help:  ## Display this help menu.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
