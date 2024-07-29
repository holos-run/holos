package token

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/token"
	"github.com/spf13/cobra"
)

// New returns a new login command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("token")
	cmd.Short = "write id token to stdout"
	cmd.Long = "Useful with curl / grpcurl -H $(holos token)"

	config := token.NewConfig()
	cmd.Flags().AddGoFlagSet(config.FlagSet())

	fs := &flag.FlagSet{}
	cmd.Flags().AddGoFlagSet(fs)
	var printClaims bool
	fs.BoolVar(&printClaims, "print-claims", false, "print id token claims")
	cmd.Flags().AddGoFlagSet(fs)

	cmd.RunE = func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		token, err := token.Get(ctx, cfg.Logger(), config)
		if err != nil {
			slog.Error("could not get token", "err", err)
			return fmt.Errorf("could not get token: %w", err)
		}

		if printClaims {
			fmt.Fprintln(cmd.OutOrStdout(), token.Pretty)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), token.Bearer)
		}

		return nil
	}

	return cmd
}
