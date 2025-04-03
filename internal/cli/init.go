package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/generate"
	"github.com/spf13/cobra"
)

// New returns a new init command.
func newInitCommand() *cobra.Command {
	cmd := command.New("init")
	cmd.Aliases = []string{"initialize", "gen", "generate"}
	cmd.Short = "initialize platforms and components"
	cmd.Args = cobra.NoArgs

	cmd.AddCommand(newInitPlatformCommand())

	return cmd
}

func newInitPlatformCommand() *cobra.Command {
	var force bool

	cmd := command.New("platform [flags] PLATFORM")
	cmd.Short = "initialize a platform from an embedded schematic"
	cmd.Long = fmt.Sprintf("Available platforms:\n\n  %s", strings.Join(generate.Platforms(), "\n  "))
	cmd.Example = "  holos init platform v1alpha5"
	cmd.Args = cobra.ExactArgs(1)

	cmd.Flags().BoolVarP(&force, "force", "", force, "force initialization")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()

		if !force {
			files, err := os.ReadDir(".")
			if err != nil {
				return errors.Wrap(err)
			}
			if len(files) > 0 {
				return errors.Format("could not initialize: directory not empty and --force=false")
			}
		}

		wd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err)
		}

		for _, name := range args {
			if err := generate.GeneratePlatform(ctx, wd, name); err != nil {
				return errors.Wrap(err)
			}
		}

		return nil
	}

	return cmd
}
