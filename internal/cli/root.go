package cli

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/holos-run/holos/version"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/platform"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cli/render"
	"github.com/holos-run/holos/internal/cli/txtar"

	cueCmd "cuelang.org/go/cmd/cue/cmd"
	cue_errors "cuelang.org/go/cue/errors"
)

//go:embed help.txt
var helpLong string

// New returns a new root *cobra.Command for command line execution.
func New(cfg *holos.Config) *cobra.Command {
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

	// Hide the help command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	rootCmd.PersistentFlags().Lookup("help").Hidden = true

	// subcommands
	rootCmd.AddCommand(render.New(cfg))
	rootCmd.AddCommand(newInitCommand())

	// Maybe not needed?
	rootCmd.AddCommand(txtar.New(cfg))

	// CUE
	rootCmd.AddCommand(newCueCmd())

	// Show
	rootCmd.AddCommand(NewShowCmd(platform.NewConfig()))

	// Compare
	rootCmd.AddCommand(NewCompareCmd())

	// Compile
	rootCmd.AddCommand(NewCompileCmd())

	return rootCmd
}

func newCueCmd() (cmd *cobra.Command) {
	// Get a handle on the cue root command fields.
	root, _ := cueCmd.New([]string{})
	// Copy the fields to our embedded command.
	cmd = command.New("cue")
	cmd.Short = root.Short
	cmd.Long = root.Long
	// Pass all arguments through to RunE.
	cmd.DisableFlagParsing = true
	cmd.Args = cobra.ArbitraryArgs

	// We do it this way so we handle errors correctly.
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cueRootCommand, _ := cueCmd.New(args)
		return cueRootCommand.Run(cmd.Root().Context())
	}
	return cmd
}

// HandleError is the top level error handler that unwraps and logs errors.
func HandleError(ctx context.Context, err error, hc *holos.Config) (exitCode int) {
	log := logger.FromContext(ctx)
	var cueErr cue_errors.Error
	var errAt *errors.ErrorAt

	if errors.As(err, &errAt) {
		loc := errAt.Source.Loc()
		err2 := errAt.Unwrap()
		log.ErrorContext(ctx, fmt.Sprintf("error at %s: %s", loc, err2), "err", err2, "loc", loc)
	} else {
		log.ErrorContext(ctx, err.Error(), "err", err)
	}

	// cue errors are bundled up as a list and refer to multiple files / lines.
	if errors.As(err, &cueErr) {
		msg := cue_errors.Details(cueErr, nil)
		if _, err := fmt.Fprint(hc.Stderr(), msg); err != nil {
			log.ErrorContext(ctx, "could not write CUE error details: "+err.Error(), "err", err)
		}
	}
	return 1
}
