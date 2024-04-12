package render

import (
	"fmt"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/internal/builder"
	"github.com/holos-run/holos/internal/logger"
	"github.com/spf13/cobra"
)

func makeRenderRunFunc(cfg *holos.Config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		if cfg.ClusterName() == "" {
			return errors.Wrap(fmt.Errorf("missing cluster name"))
		}

		ctx := cmd.Context()
		log := logger.FromContext(ctx)
		build := builder.New(builder.Entrypoints(args), builder.Cluster(cfg.ClusterName()))
		results, err := build.Run(cmd.Context())
		if err != nil {
			return errors.Wrap(err)
		}
		// TODO: Avoid accidental over-writes if to holos component instances result in
		// the same file path. Write files into a blank temporary directory, error if a
		// file exists, then move the directory into place.
		for _, result := range results {
			if result.Skip {
				continue
			}
			// API Objects
			path := result.Filename(cfg.WriteTo(), cfg.ClusterName())
			if err := result.Save(ctx, path, result.AccumulatedOutput()); err != nil {
				return errors.Wrap(err)
			}
			// Kustomization
			path = result.KustomizationFilename(cfg.WriteTo(), cfg.ClusterName())
			if err := result.Save(ctx, path, result.KsContent); err != nil {
				return errors.Wrap(err)
			}
			log.InfoContext(ctx, "rendered "+result.Name(), "status", "ok", "action", "rendered", "name", result.Name())
		}
		return nil
	}
}

// New returns the render subcommand for the root command
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("render [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "write kubernetes api objects to the filesystem"
	cmd.Flags().SortFlags = false
	cmd.Flags().AddGoFlagSet(cfg.WriteFlagSet())
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())
	cmd.RunE = makeRenderRunFunc(cfg)
	return cmd
}