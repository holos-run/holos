package initialize

import (
	"flag"
	"fmt"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/spf13/cobra"
)

// Config holds configuration parameters for initialize.
type config struct {
	schematic *string
}

// Build the shared configuration and flagset for the init subcommand.
func newConfig() (*config, *flag.FlagSet) {
	cfg := &config{}
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg.schematic = flagSet.String("schematic", "bare", "The name of the schematic being used to initialize the platform.")

	return cfg, flagSet
}

// makeInitFunc returns the internal implementation of the init cli subcommand.
func makeInitializeFunc(_ *holos.Config, cfg *config) command.RunFunc {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		log.Info("Starting Holos platform initialization...")
		fmt.Printf("Schematic: %s\n", *cfg.schematic)
		log.Info("Holos platform initialization complete.")
		return nil
	}
}

// New returns the init subcommand for the root command.
func New(hc *holos.Config) *cobra.Command {
	cmd := command.New("init")

	cfg, flagSet := newConfig()

	cmd.Short = "Initialize a new Holos platform based on a schematic."
	cmd.Flags().AddGoFlagSet(flagSet)
	cmd.RunE = makeInitializeFunc(hc, cfg)

	return cmd
}
