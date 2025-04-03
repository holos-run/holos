# tl;dr:
#
# 1. Run 'make bootstrap' to install nix and direnv
# 2. Run 'nix develop' to enter the development environment
# 3. Use 'direnv allow' or 'nix develop' to enter the development environment
# 4. Use 'direnv revoke' or 'exit' to leave the development environment
#
# This Makefile helps bootstrap a development environment with nix and direnv.
# After this is complete, you can use standard nix commands to develop and build
# the project. It is not intended to be used for anything except for bootstrapping.

.DEFAULT_GOAL := help

#-------
##@ help
#-------

# based on "https://gist.github.com/prwhite/8168133?permalink_comment_id=4260260#gistcomment-4260260"
.PHONY: help
help: ## Display this help. (Default)
	@grep -hE '^(##@|[A-Za-z0-9_ \-]*?:.*##).*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; /^##@/ {print "\n" substr($$0, 5)} /^[A-Za-z0-9_ \-]*?:.*##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

help-sort: ## Display alphabetized version of help (no section headings).
	@grep -hE '^[A-Za-z0-9_ \-]*?:.*##.*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; /^[A-Za-z0-9_ \-]*?:.*##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

HELP_TARGETS_PATTERN ?= test
help-targets: ## Print commands for all targets matching a given pattern. eval "$(make help-targets HELP_TARGETS_PATTERN=test | sed 's/\x1b\[[0-9;]*m//g')"
	@make help-sort | awk '{print $$1}' | grep '$(HELP_TARGETS_PATTERN)' | xargs -I {} printf "printf '___\n\n{}:\n\n'\nmake -n {}\nprintf '\n'\n"

# catch-all pattern rule
#
# This rule matches any targets that are not explicitly defined in this
# Makefile. It prevents 'make' from failing due to unrecognized targets, which
# is particularly useful when passing arguments or targets to sub-Makefiles. The
# '@:' command is a no-op, indicating that nothing should be done for these
# targets within this Makefile.
#
%:
	@:

#-------
##@ bootstrap
#-------

.PHONY: bootstrap
bootstrap: ## Main bootstrap target that runs all necessary setup steps
bootstrap: install-nix install-direnv
	@echo "\nBootstrap of nix and direnv complete! Please note:"
	@echo ""
	@echo "- Start a new shell session before continuing"
	@echo "- Run 'nix develop' to enter the development environment"
	@echo ""
	@echo "- If you would like to automatically activate the development environment when you enter the project directory"
	@echo "  - see https://direnv.net/docs/hook.html to add direnv to your shell"
	@echo "  - start a new shell session"
	@echo "  - 'cd' out and back into the project directory"
	@echo "  - allow direnv by running 'direnv allow'"

.PHONY: install-nix
install-nix: ## Install Nix using the Determinate Systems installer
	@echo "Installing Nix..."
	@if command -v nix >/dev/null 2>&1; then \
		echo "Nix is already installed."; \
	else \
		curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install; \
	fi

.PHONY: install-direnv
install-direnv: ## Install direnv using the official installation script
	@echo "Installing direnv..."
	@if command -v direnv >/dev/null 2>&1; then \
		echo "direnv is already installed."; \
	else \
		curl -sfL https://direnv.net/install.sh | bash; \
	fi
	@echo ""
	@echo "See https://direnv.net/docs/hook.html if you would like to add direnv to your shell"

#-------
##@ clean
#-------

.PHONY: clean
clean: ## Clean any temporary files or build artifacts
	@echo "Cleaning up..."
	@rm -rf result result-*
