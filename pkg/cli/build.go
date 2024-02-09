package cli

import (
	"fmt"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/internal/builder"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
	"strings"
)

// build is the internal implementation of the build cli command
func build(cmd *cobra.Command, args []string) error {
	build := builder.New(builder.Entrypoints(args))
	results, err := build.Run(cmd.Context())
	if err != nil {
		return err
	}
	outs := make([]string, 0, len(results))
	for _, result := range results {
		outs = append(outs, result.Content)
	}
	out := strings.Join(outs, "---\n")
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), out); err != nil {
		return wrapper.Wrap(err)
	}
	return nil
}

// newBuildCmd returns the build subcommand for the root command
func newBuildCmd(cfg *config.Config) *cobra.Command {
	cmd := newCmd("build [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "build kubernetes api objects from a directory"
	cmd.RunE = build
	return cmd
}
