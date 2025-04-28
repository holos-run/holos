package cli

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	v1alpha5 "github.com/holos-run/holos/api/core/v1alpha5"
	v1alpha6 "github.com/holos-run/holos/api/core/v1alpha6"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/compile"
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
		format: "yaml",
		cfg:    cfg,
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
	format string
	cfg    *platform.Config
}

func (s *showBuildPlans) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&s.format, "format", "yaml", "yaml or json format")
	return fs
}

func (s *showBuildPlans) Run(ctx context.Context, p *platform.Platform) error {
	components := p.Select(s.cfg.ComponentSelectors...)
	reqs := make([]compile.BuildPlanRequest, len(components))

	for idx, c := range components {
		tags, err := c.Tags()
		if err != nil {
			return errors.Wrap(err)
		}
		reqs[idx] = compile.BuildPlanRequest{
			APIVersion: "v1alpha6",
			Kind:       "BuildPlanRequest",
			Root:       p.Root(),
			Leaf:       c.Path(),
			WriteTo:    s.cfg.WriteTo,
			TempDir:    "${TMPDIR_PLACEHOLDER}",
			Tags:       tags,
		}
	}

	resp, err := compile.Compile(ctx, s.cfg.Concurrency, reqs)
	if err != nil {
		return errors.Wrap(err)
	}

	encoder, err := holos.NewEncoder(s.format, s.cfg.Stdout)
	if err != nil {
		return errors.Wrap(err)
	}

	for _, buildPlanResponse := range resp {
		var tm holos.TypeMeta
		if err := json.Unmarshal(buildPlanResponse.RawMessage, &tm); err != nil {
			return errors.Format("could not discriminate type meta: %w", err)
		}
		if tm.Kind != "BuildPlan" {
			return errors.Format("invalid kind %s: must be BuildPlan", tm.Kind)
		}

		var buildPlan any
		switch tm.APIVersion {
		case "v1alpha5":
			buildPlan = &v1alpha5.BuildPlan{}
		case "v1alpha6":
			buildPlan = &v1alpha6.BuildPlan{}
		default:
			slog.WarnContext(ctx, fmt.Sprintf("unknown BuildPlan APIVersion %s: assuming v1alpha6 schema", tm.APIVersion))
			buildPlan = &v1alpha6.BuildPlan{}
		}

		if err := json.Unmarshal(buildPlanResponse.RawMessage, buildPlan); err != nil {
			return errors.Wrap(err)
		}
		if err := encoder.Encode(buildPlan); err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}
