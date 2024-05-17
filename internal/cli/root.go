package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/holos-run/holos/version"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/server"

	"github.com/holos-run/holos/internal/cli/build"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/controller"
	"github.com/holos-run/holos/internal/cli/create"
	"github.com/holos-run/holos/internal/cli/generate"
	"github.com/holos-run/holos/internal/cli/get"
	"github.com/holos-run/holos/internal/cli/kv"
	"github.com/holos-run/holos/internal/cli/login"
	"github.com/holos-run/holos/internal/cli/logout"
	"github.com/holos-run/holos/internal/cli/preflight"
	"github.com/holos-run/holos/internal/cli/push"
	"github.com/holos-run/holos/internal/cli/register"
	"github.com/holos-run/holos/internal/cli/render"
	"github.com/holos-run/holos/internal/cli/rpc"
	"github.com/holos-run/holos/internal/cli/token"
	"github.com/holos-run/holos/internal/cli/txtar"
)

// New returns a new root *cobra.Command for command line execution.
func New(cfg *holos.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "holos",
		Short:   "holos manages a holistic integrated software development platform",
		Version: version.Version,
		Args:    cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true, // Don't complete the complete subcommand itself
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if err := cfg.Finalize(); err != nil {
				return err
			}
			log := cfg.Logger()
			c.Root().SetContext(logger.NewContext(c.Context(), log))
			// Set the default logger after flag parsing.
			slog.SetDefault(log)
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			return c.Usage()
		},
	}
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.SetOut(cfg.Stdout())
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().AddGoFlagSet(cfg.LogFlagSet())

	// subcommands
	rootCmd.AddCommand(build.New(cfg))
	rootCmd.AddCommand(render.New(cfg))
	rootCmd.AddCommand(get.New(cfg))
	rootCmd.AddCommand(create.New(cfg))
	rootCmd.AddCommand(preflight.New(cfg))
	rootCmd.AddCommand(login.New(cfg))
	rootCmd.AddCommand(logout.New(cfg))
	rootCmd.AddCommand(token.New(cfg))
	rootCmd.AddCommand(rpc.New(cfg))
	rootCmd.AddCommand(generate.New(cfg))
	rootCmd.AddCommand(register.New(cfg))
	rootCmd.AddCommand(push.New(cfg))
	rootCmd.AddCommand(newOrgCmd())

	// Maybe not needed?
	rootCmd.AddCommand(txtar.New(cfg))

	// Deprecated, remove?
	rootCmd.AddCommand(kv.New(cfg))

	// Server
	rootCmd.AddCommand(server.New(cfg))

	// Controller
	rootCmd.AddCommand(controller.New(cfg))

	return rootCmd
}

func newOrgCmd() (cmd *cobra.Command) {
	cmd = command.New("orgid")
	cmd.Short = "print the current context org id."
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		cc := holos.NewClientContext(ctx)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), cc.OrgID)
		return err
	}
	return cmd
}
