package create

import (
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/secret"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/spf13/cobra"
)

// New returns the create command for the cli
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("create")
	cmd.Short = "create resources"
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}

	// api client config
	config := client.NewConfig(cfg)

	// flags
	cmd.PersistentFlags().SortFlags = false
	// commands
	cmd.AddCommand(secret.NewCreateCmd(cfg))
	cmd.AddCommand(NewPlatform(config))
	return cmd
}

func NewPlatform(cfg *client.Config) *cobra.Command {
	cmd := command.New("platform")

	cmd.Short = "create a platform"
	cmd.Args = cobra.NoArgs

	pm := client.PlatformMutation{}
	cmd.Flags().AddGoFlagSet(pm.FlagSet())

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		client := client.New(cfg)
		pf, err := client.CreatePlatform(ctx, pm)
		if err != nil {
			return err
		}
		log := logger.FromContext(ctx)
		log.InfoContext(ctx, "created platform", "name", pf.GetName(), "id", pf.GetId(), "org", pf.GetOwner().GetOrgId())
		return nil
	}

	return cmd
}
