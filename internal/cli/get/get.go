package get

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/secret"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/slice"
)

// New returns the get command for the cli.
func New(hc *holos.Config) *cobra.Command {
	cmd := command.New("get")
	cmd.Short = "get resources"
	cmd.Aliases = []string{"list"}
	cmd.Flags().SortFlags = false
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return c.Usage()
	}
	// flags
	cmd.PersistentFlags().SortFlags = false

	// commands
	cmd.AddCommand(secret.NewGetCmd(hc))
	cmd.AddCommand(NewPlatform(hc))
	return cmd
}

func NewPlatform(hc *holos.Config) *cobra.Command {
	cmd := command.New("platform")
	cmd.Aliases = []string{"platforms"}
	cmd.Short = "rpc PlatformService.ListPlatforms"
	cmd.Long = "List platforms stored in Holos Server."
	cmd.Args = cobra.MinimumNArgs(0)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		clientContext := holos.NewClientContext(ctx)
		rpc := client.New(client.NewConfig(hc))

		msgs, err := rpc.Platforms(ctx, clientContext.OrgID)
		if err != nil {
			return err
		}

		now := time.Now()
		rows := make([][]string, 0, len(msgs))
		for _, msg := range msgs {
			name := msg.GetName()
			if len(args) > 0 && !slice.ContainsString(args, name, nil) {
				continue
			}
			age := now.Sub(msg.GetDetail().GetCreatedAt().AsTime())
			rows = append(rows, []string{
				name,
				msg.GetDisplayName(),
				holos.RoundDuration(age),
				msg.GetId(),
			})
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 4, ' ', 0)
		if len(rows) == 0 {
			return errors.New("not found")
		}
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tAGE\tID")
		for _, row := range rows {
			fmt.Fprintln(w, strings.Join(row, "\t"))
		}
		return w.Flush()
	}

	return cmd
}
