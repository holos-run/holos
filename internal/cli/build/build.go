package build

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

// makeBuildRunFunc returns the internal implementation of the build cli command
func makeBuildRunFunc(cfg *holos.Config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		build := builder.New(builder.Entrypoints(args), builder.Cluster(cfg.ClusterName()))
		results, err := build.Run(cmd.Context())
		if err != nil {
			return err
		}
		outs := make([]string, 0, len(results))
		for idx, result := range results {
			if result == nil || result.Skip {
				slog.Debug("skip result", "idx", idx, "result", result)
				continue
			}
			slog.Debug("append result", "idx", idx, "result.kind", result.Kind)
			outs = append(outs, result.AccumulatedOutput())
		}
		out := strings.Join(outs, "---\n")
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), out); err != nil {
			return errors.Wrap(err)
		}
		return nil
	}
}

// New returns the build subcommand for the root command
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("build [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "build kubernetes api objects from a directory"
	cmd.RunE = makeBuildRunFunc(cfg)
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())
	return cmd
}
