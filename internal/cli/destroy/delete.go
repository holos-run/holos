package destroy

import (
	"fmt"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

// New returns the command for the cli
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("delete")
	cmd.Aliases = []string{"destroy"}
	cmd.Short = "delete resources"
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}
	// api client config
	config := client.NewConfig(cfg)
	// flags
	cmd.PersistentFlags().SortFlags = false
	// commands
	cmd.AddCommand(NewPlatform(config))
	return cmd
}

func NewPlatform(cfg *client.Config) *cobra.Command {
	cmd := command.New("platform")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Use = "platform [flags] PLATFORM_ID [PLATFORM_ID...]"
	cmd.Short = "rpc PlatformService.DeletePlatform"

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		rpc := client.New(cfg)

		for _, platformID := range args {
			msg, err := rpc.DeletePlatform(ctx, platformID)
			if err != nil {
				return errors.Wrap(err)
			}
			platform := msg.GetPlatform()
			fmt.Fprintf(cmd.OutOrStdout(), "deleted platform %s (%s)\n", platform.GetName(), platform.GetId())
		}

		return nil
	}

	return cmd
}
