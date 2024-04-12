package login

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
	cmd := command.New("login")
	var printClaims bool

	config := token.NewConfig()
	cmd.Flags().AddGoFlagSet(config.FlagSet())

	fs := &flag.FlagSet{}
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

		claims := token.Claims()
		slog.Info("logged in as "+claims.Email, "name", claims.Name, "exp", token.Expiry, "email", claims.Email)
		if printClaims {
			fmt.Fprintln(cmd.OutOrStdout(), token.Pretty)
		}
		return nil
	}

	return cmd
}
