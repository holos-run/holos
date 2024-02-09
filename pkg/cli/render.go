package cli

import (
	"fmt"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/internal/builder"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
)

func makeRenderRunFunc(cfg *config.Config) runFunc {
	return func(cmd *cobra.Command, args []string) error {
		if cfg.ClusterName() == "" {
			return wrapper.Wrap(fmt.Errorf("missing cluster name"))
		}

		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		build := builder.New(builder.Entrypoints(args))
		results, err := build.Run(cmd.Context())
		if err != nil {
			return wrapper.Wrap(err)
		}
		for _, result := range results {
			path := result.Filename(cfg.WriteTo(), cfg.ClusterName())
			if err := result.Save(ctx, path); err != nil {
				return wrapper.Wrap(err)
			}
			log.InfoContext(ctx, "wrote", "status", "ok", "action", "save", "path", path)
		}
		return nil
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
