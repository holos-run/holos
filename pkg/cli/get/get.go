package get

import (
	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/cli/secret"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/spf13/cobra"
)

// New returns the get command for the cli.
func New(hc *holos.Config) *cobra.Command {
	cmd := command.New("get")
	cmd.Short = "get resources"
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}
	// flags
	cmd.PersistentFlags().SortFlags = false
	// commands
	cmd.AddCommand(secret.NewGetCmd(hc))
	return cmd
}
