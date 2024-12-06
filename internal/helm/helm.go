package helm

import (
	"context"
	"fmt"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

// PullChart downloads and caches a Helm chart locally. It handles both OCI and
// HTTP repositories.  Returns the path to the cached chart and any error
// encountered.  Attribution: Helm SDK Examples [Pull Action].
//
// For convenience, initialize SDK setting via CLI mechanism:
//
//	settings := cli.New()
//
// [Pull Action]: https://helm.sh/docs/sdk/examples/#pull-action
func PullChart(ctx context.Context, settings *cli.EnvSettings, chartRef, chartVersion, repoURL, destDir string) error {
	log := logger.FromContext(ctx)
	actionConfig, err := initActionConfig(ctx, settings)
	if err != nil {
		return errors.Format("failed to init action config: %w", err)
	}

	registryClient, err := newRegistryClient(settings, false)
	if err != nil {
		return errors.Format("failed to created registry client: %w", err)
	}
	actionConfig.RegistryClient = registryClient

	pullClient := action.NewPullWithOpts(action.WithConfig(actionConfig))
	pullClient.Untar = true
	pullClient.RepoURL = repoURL
	pullClient.DestDir = destDir
	pullClient.Settings = settings
	pullClient.Version = chartVersion

	result, err := pullClient.Run(chartRef)
	if err != nil {
		return errors.Format("failed to pull chart: %w", err)
	}

	log.DebugContext(ctx, fmt.Sprintf("%+v", result))

	return nil
}
