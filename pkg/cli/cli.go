package cli

import (
	"fmt"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/internal/builder"
	"github.com/holos-run/holos/pkg/version"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
)

// newCmd returns a new subcommand
func newCmd(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     name,
		Version: version.Version,
		Args:    cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		RunE: func(c *cobra.Command, args []string) error {
			return wrapper.Wrap(fmt.Errorf("could not run %v: not implemented", c.Name()))
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

// build is the internal implementation of the build cli command
func build(cmd *cobra.Command, args []string) error {
	opts := builder.Options{Entrypoints: args}
	build := builder.New(opts)
	return build.Run(cmd.Context())
}

// newBuildCmd returns the build subcommand for the root command
func newBuildCmd(cfg *config.Config) *cobra.Command {
	cmd := newCmd("build [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "build kubernetes api objects from a directory"
	cmd.RunE = build
	return cmd
}
