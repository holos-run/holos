// Package push pushes resources to the holos api server.
package push

import (
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/push"
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
	// cmd.AddCommand(NewPlatformModel(cfg))

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
		for _, name := range args {
			if err := push.PlatformForm(ctx, name); err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	}

	return cmd
}
