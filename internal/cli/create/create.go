package create

import (
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/secret"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

// New returns the create command for the cli
func New(hc *holos.Config) *cobra.Command {
	cmd := command.New("create")
	cmd.Short = "create resources"
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}
	// flags
	cmd.PersistentFlags().SortFlags = false
	// commands
	cmd.AddCommand(secret.NewCreateCmd(hc))
	return cmd
}
