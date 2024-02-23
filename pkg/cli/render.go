package cli

import (
	"fmt"
	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/internal/builder"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
)

func makeRenderRunFunc(cfg *config.Config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		if cfg.ClusterName() == "" {
			return wrapper.Wrap(fmt.Errorf("missing cluster name"))
		}

		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		build := builder.New(builder.Entrypoints(args), builder.Cluster(cfg.ClusterName()))
		results, err := build.Run(cmd.Context())
		if err != nil {
			return wrapper.Wrap(err)
		}
		// TODO: Avoid accidental over-writes if to holos component instances result in
		// the same file path. Write files into a blank temporary directory, error if a
		// file exists, then move the directory into place.
		for _, result := range results {
			// API Objects
			path := result.Filename(cfg.WriteTo(), cfg.ClusterName())
			if err := result.Save(ctx, path, result.Content); err != nil {
				return wrapper.Wrap(err)
			}
			// Kustomization
			path = result.KustomizationFilename(cfg.WriteTo(), cfg.ClusterName())
			if err := result.Save(ctx, path, result.KsContent); err != nil {
				return wrapper.Wrap(err)
			}
			log.InfoContext(ctx, "rendered "+result.Name(), "status", "ok", "action", "rendered", "name", result.Name())
		}
		return nil
	}
}

// newRenderCmd returns the render subcommand for the root command
func newRenderCmd(cfg *config.Config) *cobra.Command {
	cmd := command.New("render [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "write kubernetes api objects to the filesystem"
	cmd.Flags().SortFlags = false
	cmd.Flags().AddGoFlagSet(cfg.WriteFlagSet())
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())
	cmd.RunE = makeRenderRunFunc(cfg)
	return cmd
}
