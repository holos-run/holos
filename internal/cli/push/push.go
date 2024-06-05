// Package push pushes resources to the holos api server.
package push

import (
	"fmt"
	"log/slog"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/push"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("push")
	cmd.Short = "push resources to holos server"
	cmd.Args = cobra.NoArgs

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	cmd.AddCommand(NewPlatform(config))

	return cmd
}

func NewPlatform(cfg *client.Config) *cobra.Command {
	cmd := command.New("platform")

	cmd.Short = "push platform resources to holos server"
	cmd.Args = cobra.NoArgs

	cmd.AddCommand(NewPlatformForm(cfg))
	cmd.AddCommand(NewPlatformModel(cfg))

	return cmd
}

func NewPlatformForm(cfg *client.Config) *cobra.Command {
	cmd := command.New("form")
	cmd.Short = "push platform form to holos server"
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
			p, err := client.LoadPlatform(ctx, name)
			if err != nil {
				return errors.Wrap(err)
			}
			// Build the form from the cue code.
			form, err := push.PlatformForm(ctx, name)
			if err != nil {
				return errors.Wrap(err)
			}
			// Make the rpc call to update the platform form.
			if err := rpc.UpdateForm(ctx, p.GetId(), form); err != nil {
				return errors.Wrap(err)
			}
			slog.Default().InfoContext(ctx, fmt.Sprintf("pushed: %s/ui/platform/%s", cfg.Client().Server(), p.GetId()))
		}
		return nil
	}

	return cmd
}

func NewPlatformModel(cfg *client.Config) *cobra.Command {
	cmd := command.New("model")
	cmd.Short = "push platform model to holos server"
	cmd.Args = cobra.MinimumNArgs(1)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		if ctx == nil {
			return errors.Wrap(errors.New("cannot execute: no context"))
		}
		ctx = logger.NewContext(ctx, logger.FromContext(ctx).With("server", cfg.Client().Server()))
		rpc := client.New(cfg)
		for _, name := range args {
			// Get the platform config for the platform id.
			p, err := client.LoadPlatformConfig(ctx, name)
			if err != nil {
				return errors.Wrap(err)
			}

			// Make the rpc call to update the platform form.
			if err := rpc.UpdatePlatformModel(ctx, p.PlatformId, p.PlatformModel); err != nil {
				return errors.Wrap(err)
			}
			slog.Default().InfoContext(ctx, fmt.Sprintf("pushed: %s/ui/platform/%s", cfg.Client().Server(), p.PlatformId))
		}
		return nil
	}

	return cmd
}
