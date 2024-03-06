package preflight

import (
	"flag"

	"github.com/spf13/cobra"

	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/holos-run/holos/pkg/logger"
)

// Config holds configuration parameters for preflight checks.
type config struct {
	githubInstance *string
}

// Build the shared configuration and flagset for the preflight command.
func newConfig() (*config, *flag.FlagSet) {
	cfg := &config{}
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg.githubInstance = flagSet.String("github-instance", "github.com", "Address of the GitHub instance you want to use")

	return cfg, flagSet
}

// New returns the preflight command for the root command.
func New(hc *holos.Config) *cobra.Command {
	cfg, flagSet := newConfig()

	cmd := command.New("preflight")
	cmd.Short = "Run preflight checks to ensure you're ready to use Holos"
	cmd.Flags().AddGoFlagSet(flagSet)
	cmd.RunE = makePreflightRunFunc(hc, cfg)

	return cmd
}

// makePreflightRunFunc returns the internal implementation of the preflight cli command.
func makePreflightRunFunc(_ *holos.Config, cfg *config) command.RunFunc {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		log.Info("Starting preflight checks")

		// GitHub checks
		if err := RunGhChecks(ctx, cfg); err != nil {
			return err
		}
		// Other checks can be added here

		log.Info("Preflight checks complete. Ready to use Holos 🚀")
		return nil
	}
}
