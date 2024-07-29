// Package register provides user registration via the command line.
package register

import (
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/register"
	"github.com/spf13/cobra"
)

// New returns a new register command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("register")
	cmd.Short = "rpc UserService.RegisterUser"
	cmd.Long = "register with holos server"
	cmd.Args = cobra.NoArgs

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	cmd.AddCommand(NewUser(config))

	return cmd
}

// NewUser returns a command to register a user with holos server.
func NewUser(cfg *client.Config) *cobra.Command {
	cmd := command.New("user")
	cmd.Short = "user registration workflow"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		return register.User(ctx, cfg)
	}
	return cmd
}
