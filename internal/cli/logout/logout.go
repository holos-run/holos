package logout

import (
	"fmt"
	"os"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/token"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("logout")
	cmd.Hidden = !feature.Flag(holos.ServerFeature)
	cmd.Short = "log out by deleting cached credentials"
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if err := os.RemoveAll(token.CacheDir); err != nil {
			return errors.Wrap(fmt.Errorf("could not logout: %w", err))
		}
		cfg.Logger().Info("logged out: removed " + token.CacheDir)
		return nil
	}
	return cmd
}
