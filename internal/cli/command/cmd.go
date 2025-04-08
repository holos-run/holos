package command

import (
	"github.com/holos-run/holos/version"
	"github.com/spf13/cobra"
)

// RunFunc is a cobra.Command RunE function.
type RunFunc func(c *cobra.Command, args []string) error

// New returns a new subcommand
func New(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     name,
		Short:   name,
		Version: version.GetVersion(),
		Args:    cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().SortFlags = false
	return cmd
}
