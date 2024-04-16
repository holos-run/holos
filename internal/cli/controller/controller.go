// Package controller integrates Choria Machine Room into Holos for cluster management.
package controller

import (
	"context"
	"fmt"

	mr "github.com/choria-io/machine-room"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/version"
	"github.com/spf13/cobra"
)

var (
	// SigningKey is the public key from choria jwt keys machine-signer.seed machine-signer.public, refer to gen-machine-signer.
	SigningKey = "2a136e3875f4375968ae8e8d400ba24864d3ed7c4109675f357d32cc3ca1d5a7"
)

func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("controller")
	cmd.Args = cobra.ArbitraryArgs
	cmd.DisableFlagParsing = true
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if SigningKey == "" {
			return errors.Wrap(fmt.Errorf("could not run: controller.SigningKey not set from build variables"))
		}

		ctx := c.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		app, err := mr.New(mr.Options{
			Name:              "controller",
			Contact:           "jeff@openinfrastructure.co",
			Version:           version.Version,
			Help:              "Holos Controller",
			MachineSigningKey: SigningKey,
			Args:              args,
		})
		if err != nil {
			return errors.Wrap(fmt.Errorf("could not make machine room app: %w", err))
		}

		return app.Run(ctx)
	}
	return cmd
}
