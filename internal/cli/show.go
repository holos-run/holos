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

func NewShowCmd(cfg platform.Config) (cmd *cobra.Command) {
	cmd = command.New("show")
	cmd.Short = "show a platform or build plans"

	spf := &showPlatform{
		Format: "yaml",
		Out:    cfg.Stdout,
	}
	platformCmd := platform.NewCommand(cfg, spf.Run)
	platformCmd.Flags().AddFlagSet(spf.flagSet())
	cmd.AddCommand(platformCmd)

	sbp := &showBuildPlans{
		Format: "yaml",
		Out:    cfg.Stdout,
	}
	buildPlanCmd := platform.NewCommand(platform.NewConfig(), sbp.Run)
	buildPlanCmd.Use = "buildplans"
	buildPlanCmd.Short = "show buildplans"
	buildPlanCmd.Long = longShowBuildPlansHelp
	buildPlanCmd.Aliases = []string{"buildplan", "components", "component"}
	buildPlanCmd.Flags().AddFlagSet(sbp.flagSet())
	cmd.AddCommand(buildPlanCmd)
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
	Format    string
	Out       io.Writer
	Selectors holos.Selectors
}

func (s *showBuildPlans) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&s.Format, "format", "yaml", "yaml or json format")
	fs.VarP(&s.Selectors, "selector", "l", "label selector (e.g. label==string,label!=string)")
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
			// Export the build plan using the ordered encoder.
			return errors.Wrap(bp.Export(idx, encoder))
		},
		ComponentSelectors: s.Selectors,
	}

	return errors.Wrap(p.Build(ctx, opts))
}
