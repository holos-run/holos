// Package pull pulls resources from the PlatformService and caches them in the
// local filesystem.
package pull

import (
	"fmt"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("pull")
	cmd.Hidden = !feature.Flag(holos.ServerFeature)
	cmd.Short = "pull resources from holos server"
	cmd.Args = cobra.NoArgs

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	cmd.AddCommand(NewPlatform(config))

	return cmd
}

func NewPlatform(cfg *client.Config) *cobra.Command {
	cmd := command.New("platform")

	cmd.Short = "pull platform resources"
	cmd.Args = cobra.NoArgs

	cmd.AddCommand(NewPlatformConfig(cfg))

	return cmd
}

func NewPlatformConfig(cfg *client.Config) *cobra.Command {
	cmd := command.New("model")
	cmd.Aliases = []string{"config"}
	cmd.Short = "pull platform model"
	cmd.Args = cobra.MinimumNArgs(1)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		if ctx == nil {
			return errors.Wrap(errors.New("cannot execute: no context"))
		}
		ctx = logger.NewContext(ctx, logger.FromContext(ctx).With("server", cfg.Client().Server()))
		rpc := client.New(cfg)
		for _, name := range args {
			// Get the platform metadata for the platform id.
			pmd, err := client.LoadPlatformMetadata(ctx, name)
			if err != nil {
				return errors.Wrap(err)
			}
			log := logger.FromContext(ctx).With("platform_name", pmd.GetName(), "platform_id", pmd.GetId())
			// Get the platform model
			model, err := rpc.PlatformModel(ctx, pmd.GetId())
			if err != nil {
				return errors.Wrap(err)
			}
			log.DebugContext(ctx, "pulled platform "+pmd.GetName()+" model")
			// Build the PlatformConfig
			pc := &object.PlatformConfig{PlatformModel: model}
			// Save the PlatformConfig
			path, err := client.SavePlatformConfig(ctx, name, pc)
			if err != nil {
				return errors.Wrap(err)
			}
			log.DebugContext(ctx, fmt.Sprintf("wrote: %s for platform %s", path, pmd.GetName()), "path", path)
		}
		return nil
	}

	return cmd
}
