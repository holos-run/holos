package generate

import (
	"fmt"
	"strings"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/generate"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

// New returns a new generate command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("generate")
	cmd.Aliases = []string{"gen"}
	cmd.Short = "generate local resources"
	cmd.Args = cobra.NoArgs

	cmd.AddCommand(NewPlatform(cfg))

	return cmd
}

func NewPlatform(cfg *holos.Config) *cobra.Command {
	cmd := command.New("platform")
	cmd.Use = "platform [flags] PLATFORM"
	cmd.Short = "generate a platform from an embedded schematic"
	cmd.Long = fmt.Sprintf("Embedded platforms available to generate:\n\n  %s", strings.Join(generate.Platforms(), "\n  "))
	cmd.Args = cobra.ExactArgs(1)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		clientContext := holos.NewClientContext(ctx)
		client := client.New(client.NewConfig(cfg))

		for _, name := range args {
			if err := generate.GeneratePlatform(ctx, client, clientContext.OrgID, name); err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	}

	return cmd
}
