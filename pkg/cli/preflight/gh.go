package preflight

import (
	"context"
	"fmt"
	"strings"

	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/util"
	"github.com/holos-run/holos/pkg/wrapper"
)

type ghAuthStatusResponse string

// RunGhChecks runs all the preflight checks related to GitHub.
func RunGhChecks(ctx context.Context, cfg *config) error {
	if err := cliIsInstalled(ctx); err != nil {
		return err
	}

	if err := cliIsAuthed(ctx, cfg); err != nil {
		return err
	}

	return nil
}

// cliIsInstalled checks if the GitHub CLI is installed.
func cliIsInstalled(ctx context.Context) error {
	log := logger.FromContext(ctx)

	version, err := getGhVersion(ctx)
	if err != nil {
		log.WarnContext(ctx, "GitHub CLI (gh) not installed or not in PATH.")
		return guideToInstallGh(ctx)
	}

	log.InfoContext(ctx, "GitHub CLI found", "gh_version", version)
	return nil
}

// cliIsAuthed checks if 'gh' is authenticated. If not, 'gh auth login' is run then cliIsAuthed is called again.
func cliIsAuthed(ctx context.Context, cfg *config) error {
	log := logger.FromContext(ctx)

	status, err := ghAuthStatus(ctx, cfg)

	if err != nil || !ghIsAuthenticated(status, cfg) {
		log.WarnContext(ctx, "GitHub CLI not authenticated to "+*cfg.githubInstance)
		if err := authenticateGh(ctx, cfg); err != nil {
			return wrapper.Wrap(fmt.Errorf("failed to authenticate the GitHub CLI to %v: %w", cfg.githubInstance, err))
		}

		// Re-run this check now that gh should be authenticated.
		err := cliIsAuthed(ctx, cfg)
		return err
	}

	log.InfoContext(ctx, "GitHub CLI is authenticated to "+*cfg.githubInstance)

	if !ghTokenAllowsRepoCreation(status) {
		return wrapper.Wrap(fmt.Errorf("GitHub token does not have the necessary scopes to create a repository"))
	}

	log.InfoContext(ctx, "GitHub token is able to create a repository")
	return nil
}

// ghAuthStatus runs 'gh auth status' and returns the result.
func ghAuthStatus(ctx context.Context, cfg *config) (ghAuthStatusResponse, error) {
	log := logger.FromContext(ctx)
	out, err := util.RunCmd(ctx, "gh", "auth", "status", "--hostname="+*cfg.githubInstance)

	var status ghAuthStatusResponse
	if err != nil {
		status = ghAuthStatusResponse(out.Stderr.String())
	} else {
		status = ghAuthStatusResponse(out.Stdout.String())
	}
	log.DebugContext(ctx, "gh auth status", "gh_auth_status", status)

	return status, err
}

// getGhVersion retrieves the version of 'gh'.
func getGhVersion(ctx context.Context) (string, error) {
	out, err := util.RunCmd(ctx, "gh", "--version")
	if err != nil {
		return "", err
	}
	return strings.Split(out.Stdout.String(), "\n")[0], nil
}

// guideToInstallGh guides the user towards installing the GitHub CLI.
func guideToInstallGh(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.WarnContext(ctx, "The GitHub CLI is required to set up Holos. To install it, follow the instructions at: https://github.com/cli/cli#installation")
	return wrapper.Wrap(fmt.Errorf("GitHub CLI is not installed"))
}

// authenticateGh runs 'gh auth login' to authenticate the GitHub CLI.
func authenticateGh(ctx context.Context, cfg *config) error {
	log := logger.FromContext(ctx)
	log.InfoContext(ctx, "Authenticating GitHub CLI with 'gh auth login --hostname="+*cfg.githubInstance+"'. Please follow the prompts.")

	err := util.RunInteractiveCmd(ctx, "gh", "auth", "login", "--hostname="+*cfg.githubInstance)
	if err != nil {
		log.ErrorContext(ctx, "Failed to authenticate GitHub CLI")
		return wrapper.Wrap(fmt.Errorf("failed to authenticate GitHub CLI: %w", err))
	}

	log.InfoContext(ctx, "GitHub CLI has been authenticated.")
	return nil
}

// ghIsAuthenticated checks if the GitHub CLI is authenticated and logged in to githubInstance.
func ghIsAuthenticated(status ghAuthStatusResponse, cfg *config) bool {
	return strings.Contains(string(status), "Logged in to "+*cfg.githubInstance)
}

// ghTokenAllowsRepoCreation validates that the GitHub CLI is authenticated
// with a token that allows repository creation. This is a naive implementation
// that just checks the output of 'gh auth status' for the presence of the
// 'repo' scope. Note that the 'repo' scope is sufficient to create a secret in
// a repository, so this check also covers that.
// Example token scope line: "- Token scopes: 'admin:public_key', 'gist', 'read:org', 'repo'"
func ghTokenAllowsRepoCreation(status ghAuthStatusResponse) bool {
	return strings.Contains(string(status), "'repo'")
}
