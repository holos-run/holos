package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/generate"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

// New returns a new init command.
func newInitCommand(feature holos.Flagger) *cobra.Command {
	cmd := command.New("init")
	cmd.Aliases = []string{"initialize", "gen", "generate"}
	cmd.Short = "initialize platforms and components"
	cmd.Args = cobra.NoArgs

	cmd.AddCommand(newInitPlatformCommand())
	cmd.AddCommand(newInitComponentCommand(feature))

	return cmd
}

func newInitPlatformCommand() *cobra.Command {
	var force bool

	cmd := command.New("platform [flags] PLATFORM")
	cmd.Short = "initialize a platform from an embedded schematic"
	cmd.Long = fmt.Sprintf("Available platforms:\n\n  %s", strings.Join(generate.Platforms(), "\n  "))
	cmd.Example = "  holos init platform v1alpha5"
	cmd.Args = cobra.ExactArgs(1)

	cmd.Flags().BoolVarP(&force, "force", "", force, "force initialization")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()

		if !force {
			files, err := os.ReadDir(".")
			if err != nil {
				return errors.Wrap(err)
			}
			if len(files) > 0 {
				return errors.Format("could not initialize: directory not empty and --force=false")
			}
		}

		for _, name := range args {
			if err := generate.GeneratePlatform(ctx, name); err != nil {
				return errors.Wrap(err)
			}
		}

		return nil
	}

	return cmd
}

// newInitComponentCommand returns a command to generate a holos component
func newInitComponentCommand(feature holos.Flagger) *cobra.Command {
	cmd := command.New("component")
	cmd.Short = "initialize a component from an embedded schematic"
	cmd.Hidden = !feature.Flag(holos.GenerateComponentFeature)

	for _, name := range generate.Components("v1alpha3") {
		cmd.AddCommand(makeSchematicCommand("v1alpha3", name))
	}

	return cmd
}

func makeSchematicCommand(kind, name string) *cobra.Command {
	cmd := command.New(name)
	cfg, err := generate.NewSchematic(filepath.Join("components", kind), name)
	if err != nil {
		slog.Error("could not get schematic", "err", err)
		return nil
	}
	cmd.Short = cfg.Short
	cmd.Long = cfg.Long
	cmd.Args = cobra.NoArgs
	cmd.Flags().AddGoFlagSet(cfg.FlagSet())

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		if err := generate.GenerateComponent(ctx, kind, name, cfg); err != nil {
			return errors.Wrap(err)
		}
		return nil
	}

	return cmd
}
