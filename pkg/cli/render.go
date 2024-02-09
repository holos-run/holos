package cli

import (
	"fmt"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/internal/builder"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
	"strings"
)

func makeRenderRunFunc(cfg *config.Config) runFunc {
	return func(cmd *cobra.Command, args []string) error {
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
		return wrapper.Wrap(fmt.Errorf("write the output to: %+v", cfg.WriteTo()))
	}
}

// newRenderCmd returns the render subcommand for the root command
func newRenderCmd(cfg *config.Config) *cobra.Command {
	cmd := newCmd("render [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "write kubernetes api objects to the filesystem"
	cmd.Flags().SortFlags = false
	cmd.Flags().AddGoFlagSet(cfg.WriteFlagSet())
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())
	cmd.RunE = makeRenderRunFunc(cfg)
	return cmd
}
