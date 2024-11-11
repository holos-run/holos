package cli

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/holos-run/holos/version"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/server"

	"github.com/holos-run/holos/internal/cli/build"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/create"
	"github.com/holos-run/holos/internal/cli/destroy"
	"github.com/holos-run/holos/internal/cli/get"
	"github.com/holos-run/holos/internal/cli/kv"
	"github.com/holos-run/holos/internal/cli/login"
	"github.com/holos-run/holos/internal/cli/logout"
	"github.com/holos-run/holos/internal/cli/preflight"
	"github.com/holos-run/holos/internal/cli/pull"
	"github.com/holos-run/holos/internal/cli/push"
	"github.com/holos-run/holos/internal/cli/register"
	"github.com/holos-run/holos/internal/cli/render"
	"github.com/holos-run/holos/internal/cli/token"
	"github.com/holos-run/holos/internal/cli/txtar"

	cue "cuelang.org/go/cmd/cue/cmd"
)

//go:embed help.txt
var helpLong string

// New returns a new root *cobra.Command for command line execution.
func New(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "holos",
		Short:   "holos manages a holistic integrated software development platform",
		Long:    helpLong,
		Version: version.GetVersion(),
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
	rootCmd.AddCommand(build.New(cfg, feature))
	rootCmd.AddCommand(render.New(cfg, feature))
	rootCmd.AddCommand(get.New(cfg, feature))
	rootCmd.AddCommand(create.New(cfg, feature))
	rootCmd.AddCommand(destroy.New(cfg, feature))
	rootCmd.AddCommand(preflight.New(cfg, feature))
	rootCmd.AddCommand(login.New(cfg, feature))
	rootCmd.AddCommand(logout.New(cfg, feature))
	rootCmd.AddCommand(token.New(cfg, feature))
	rootCmd.AddCommand(newInitCommand(cfg, feature))
	rootCmd.AddCommand(register.New(cfg, feature))
	rootCmd.AddCommand(pull.New(cfg, feature))
	rootCmd.AddCommand(push.New(cfg, feature))
	rootCmd.AddCommand(newOrgCmd(feature))

	// Maybe not needed?
	rootCmd.AddCommand(txtar.New(cfg))

	// Deprecated, remove?
	rootCmd.AddCommand(kv.New(cfg, feature))

	// Server
	rootCmd.AddCommand(server.New(cfg, feature))

	// CUE
	rootCmd.AddCommand(newCueCmd())

	return rootCmd
}

func newOrgCmd(feature holos.Flagger) (cmd *cobra.Command) {
	cmd = command.New("orgid")
	cmd.Short = "print the current context org id."
	cmd.Hidden = !feature.Flag(holos.ServerFeature)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		cc := holos.NewClientContext(ctx)
		_, err := fmt.Fprintln(cmd.OutOrStdout(), cc.OrgID)
		return err
	}
	return cmd
}

func newCueCmd() (cmd *cobra.Command) {
	cueCmd, _ := cue.New(os.Args[2:])
	cmd = cueCmd.Command
	return
}
