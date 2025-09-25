package helm

import (
	"context"
	"fmt"
	"net/url"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
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
func PullChart(ctx context.Context, settings *cli.EnvSettings, chartRef, chartVersion, repoURL, destDir, username, password string) error {
	log := logger.FromContext(ctx)
	actionConfig, err := initActionConfig(ctx, settings)
	if err != nil {
		return errors.Format("failed to init action config: %w", err)
	}

	registryClient, err := newDefaultRegistryClient(settings, false)
	if err != nil {
		return errors.Format("failed to created registry client: %w", err)
	}
	actionConfig.RegistryClient = registryClient

	chartRefURL, err := url.Parse(chartRef)
	if err != nil {
		return errors.Format("Failed to parse the Chart: %w", err)
	}

	// If the chart been pulled is an OCI chart, the repo authentication has to be done ahead of the pull.
	if chartRefURL.Scheme == "oci" && username != "" && password != "" {
		loginOption := registry.LoginOptBasicAuth(username, password)
		err = registryClient.Login(chartRefURL.Host, loginOption)
		if err != nil {
			return errors.Format("failed to login to registry: %w", err)
		}
	}

	pullClient := action.NewPullWithOpts(action.WithConfig(actionConfig))
	pullClient.Untar = true
	pullClient.RepoURL = repoURL
	pullClient.DestDir = destDir
	pullClient.Settings = settings
	pullClient.Version = chartVersion
	pullClient.Username = username
	pullClient.Password = password

	result, err := pullClient.Run(chartRef)
	if err != nil {
		return errors.Format("failed to pull chart: %w", err)
	}

	log.DebugContext(ctx, fmt.Sprintf("%+v", result))

	return nil
}
