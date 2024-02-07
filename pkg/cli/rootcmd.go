package cli

import (
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/version"
	"github.com/spf13/cobra"
	"log/slog"
)

// New returns a new root *cobra.Command for command line execution.
func New(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "holos",
		Short:   "holos manages a holistic integrated software development platform",
		Version: version.Version,
		Args:    cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true, // Don't complete the complete subcommand itself
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if err := cfg.Finalize(); err != nil {
				return err
			}
			log := cfg.Logger()
			// Set the configured logger in the context.
			c.SetContext(logger.NewContext(c.Context(), log))
			// Set the default logger after flag parsing.
			slog.SetDefault(log)
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			cfg.Logger().InfoContext(c.Context(), "hello")
			return nil
		},
	}
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().AddGoFlagSet(cfg.LogFlagSet())

	// build subcommand
	rootCmd.AddCommand(newBuildCmd(cfg))

	return rootCmd
}
