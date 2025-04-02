package cli

import (
	"context"
	_ "embed"
	"io"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/component"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/platform"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//go:embed long-show-buildplans.txt
var longShowBuildPlansHelp string

func NewShowCmd(cfg *platform.Config) (cmd *cobra.Command) {
	cmd = command.New("show")
	cmd.Short = "show a platform or build plans"

	spf := &showPlatform{
		Format: "yaml",
		Out:    cfg.Stdout,
	}
	spCmd := platform.NewCommand(cfg, spf.Run)
	spCmd.Flags().AddFlagSet(cfg.FlagSetTags())
	spCmd.Flags().AddFlagSet(spf.flagSet())
	cmd.AddCommand(spCmd)

	sbp := &showBuildPlans{
		Format: "yaml",
		Out:    cfg.Stdout,
	}
	sbCmd := platform.NewCommand(cfg, sbp.Run)
	sbCmd.Use = "buildplans"
	sbCmd.Short = "show buildplans"
	sbCmd.Long = longShowBuildPlansHelp
	sbCmd.Aliases = []string{"buildplan", "components", "component"}
	sbCmd.Flags().AddFlagSet(cfg.FlagSet())
	sbCmd.Flags().AddFlagSet(sbp.flagSet())
	cmd.AddCommand(sbCmd)
	return cmd
}

type showPlatform struct {
	Format string
	Out    io.Writer
}

func (s *showPlatform) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&s.Format, "format", "yaml", "yaml or json format")
	return fs
}

func (s *showPlatform) Run(ctx context.Context, p *platform.Platform) error {
	encoder, err := holos.NewEncoder(s.Format, s.Out)
	if err != nil {
		return errors.Wrap(err)
	}
	defer encoder.Close()
	return errors.Wrap(p.Export(encoder))
}

type showBuildPlans struct {
	Format string
	Out    io.Writer
}

func (s *showBuildPlans) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&s.Format, "format", "yaml", "yaml or json format")
	return fs
}

func (s *showBuildPlans) Run(ctx context.Context, p *platform.Platform) error {
	encoder, err := holos.NewSequentialEncoder(s.Format, s.Out)
	if err != nil {
		return errors.Wrap(err)
	}
	defer encoder.Close()

	opts := platform.BuildOpts{
		PerComponentFunc: func(ctx context.Context, idx int, pc holos.Component) error {
			c := component.New(p.Root(), pc.Path(), component.NewConfig())
			tm, err := c.TypeMeta()
			if err != nil {
				return errors.Wrap(err)
			}
			opts := holos.NewBuildOpts(pc.Path())
			opts.BuildContext.TempDir = "${TMPDIR_PLACEHOLDER}"

			// TODO(jjm): refactor into [holos.NewBuildOpts] as functional options.
			// Component name, label, annotations passed via tags to cue.
			tags, err := pc.Tags()
			if err != nil {
				return errors.Wrap(err)
			}
			opts.Tags = tags

			bp, err := c.BuildPlan(tm, opts)
			if err != nil {
				return errors.Wrap(err)
			}
			// Export the build plan using the sequential encoder.
			return errors.Wrap(bp.Export(idx, encoder))
		},
	}

	return errors.Wrap(p.Build(ctx, opts))
}
