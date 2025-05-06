package cli

import (
	_ "embed"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/compile"
	"github.com/holos-run/holos/internal/errors"
	"github.com/spf13/cobra"
)

//go:embed compile.txt
var compileLong string

// NewCompileCmd returns a new compile command.
func NewCompileCmd() *cobra.Command {
	cmd := command.New("compile")
	cmd.Short = "Compile Components (stdin) to BuildPlans (stdout) using CUE"
	cmd.Long = compileLong
	cmd.Args = cobra.NoArgs
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		c := compile.New()
		ctx := cmd.Root().Context()
		return errors.Wrap(c.Run(ctx))
	}
	return cmd
}
