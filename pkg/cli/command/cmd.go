package command

import (
	"fmt"
	"github.com/holos-run/holos/pkg/version"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
)

// RunFunc is a cobra.Command RunE function.
type RunFunc func(c *cobra.Command, args []string) error

// New returns a new subcommand
func New(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     name,
		Version: version.Version,
		Args:    cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		RunE: func(c *cobra.Command, args []string) error {
			return wrapper.Wrap(fmt.Errorf("could not run %v: not implemented", c.Name()))
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

// EnsureNewline adds a trailing newline if not already there.
func EnsureNewline(b []byte) []byte {
	if len(b) > 0 && b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	return b
}
