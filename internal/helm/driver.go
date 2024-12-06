package helm

import (
	"context"
	"fmt"
	"os"

	"github.com/holos-run/holos/internal/server/middleware/logger"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
)

// https://helm.sh/docs/sdk/examples/#driver
var helmDriver string = os.Getenv("HELM_DRIVER")

func initActionConfig(ctx context.Context, settings *cli.EnvSettings) (*action.Configuration, error) {
	return initActionConfigList(ctx, settings, false)
}

func initActionConfigList(ctx context.Context, settings *cli.EnvSettings, allNamespaces bool) (*action.Configuration, error) {
	log := logger.FromContext(ctx)
	actionConfig := new(action.Configuration)

	namespace := func() string {
		// For list action, you can pass an empty string instead of settings.Namespace() to list
		// all namespaces
		if allNamespaces {
			return ""
		}
		return settings.Namespace()
	}()

	debug := func(format string, a ...any) {
		log.DebugContext(ctx, fmt.Sprintf(format, a...))
	}

	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		namespace,
		helmDriver,
		debug); err != nil {
		return nil, err
	}

	return actionConfig, nil
}

func newDefaultRegistryClient(settings *cli.EnvSettings, plainHTTP bool) (*registry.Client, error) {
	opts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}
	if plainHTTP {
		opts = append(opts, registry.ClientOptPlainHTTP())
	}

	// Create a new registry client
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}
