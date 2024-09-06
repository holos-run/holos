package generate

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/generate"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

// New returns a new generate command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("generate")
	cmd.Aliases = []string{"gen"}
	cmd.Short = "generate local resources"
	cmd.Args = cobra.NoArgs

	cmd.AddCommand(NewPlatform(cfg))
	cmd.AddCommand(NewComponent())

	return cmd
}

func NewPlatform(cfg *holos.Config) *cobra.Command {
	cmd := command.New("platform [flags] PLATFORM")
	cmd.Short = "generate a platform from an embedded schematic"
	cmd.Long = fmt.Sprintf("Embedded platforms available to generate:\n\n  %s", strings.Join(generate.Platforms(), "\n  "))
	cmd.Example = "  holos generate platform k3d"
	cmd.Args = cobra.ExactArgs(1)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()

		for _, name := range args {
			if err := generate.GeneratePlatform(ctx, name); err != nil {
				return errors.Wrap(err)
			}
		}

		return nil
	}

	return cmd
}

// NewComponent returns a command to generate a holos component
func NewComponent() *cobra.Command {
	cmd := command.New("component")
	cmd.Short = "generate a component from an embedded schematic"

	cmd.AddCommand(NewCueComponent())
	cmd.AddCommand(NewHelmComponent())

	return cmd
}

func NewHelmComponent() *cobra.Command {
	cmd := command.New("helm")
	cmd.Short = "generate a helm component from a schematic"

	for _, name := range generate.HelmComponents() {
		cmd.AddCommand(makeSchematicCommand("helm", name))
	}

	return cmd
}

func NewCueComponent() *cobra.Command {
	cmd := command.New("cue")
	cmd.Short = "generate a cue component from a schematic"

	for _, name := range generate.CueComponents() {
		cmd.AddCommand(makeSchematicCommand("cue", name))
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
