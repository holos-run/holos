package cli

import (
	"context"
	_ "embed"
	"runtime"

	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
)

//go:embed long-show-buildplans.txt
var longShowBuildPlansHelp string

func newShowCmd() (cmd *cobra.Command) {
	cmd = command.New("show")
	cmd.Short = "show a platform or build plans"
	cmd.AddCommand(newShowPlatformCmd())
	cmd.AddCommand(newShowBuildPlanCmd())
	return cmd
}

func newShowPlatformCmd() (cmd *cobra.Command) {
	cmd = command.New("platform")
	cmd.Short = "show a platform"
	cmd.Args = cobra.NoArgs

	var platform string
	cmd.Flags().StringVar(&platform, "platform", "./platform", "platform directory path")
	var format string
	cmd.Flags().StringVar(&format, "format", "yaml", "yaml or json format")
	tagMap := make(holos.TagMap)
	cmd.Flags().VarP(&tagMap, "inject", "t", "set the value of a cue @tag field from a key=value pair")

	cmd.RunE = func(c *cobra.Command, args []string) (err error) {
		inst, err := builder.LoadInstance(platform, tagMap.Tags())
		if err != nil {
			return errors.Wrap(err)
		}

		encoder, err := holos.NewEncoder(format, cmd.OutOrStdout())
		if err != nil {
			return errors.Wrap(err)
		}
		defer encoder.Close()

		if err := inst.Export(encoder); err != nil {
			return errors.Wrap(err)
		}
		return nil
	}
	return cmd
}

func newShowBuildPlanCmd() (cmd *cobra.Command) {
	cmd = command.New("buildplans")
	cmd.Aliases = []string{"buildplan", "components", "component"}
	cmd.Short = "show buildplans"
	cmd.Long = longShowBuildPlansHelp
	cmd.Args = cobra.MinimumNArgs(0)

	var platform string
	cmd.Flags().StringVar(&platform, "platform", "./platform", "platform directory path")
	var format string
	cmd.Flags().StringVar(&format, "format", "yaml", "yaml or json format")
	var selector holos.Selector
	cmd.Flags().VarP(&selector, "selector", "l", "label selector (e.g. label==string,label!=string)")
	tagMap := make(holos.TagMap)
	cmd.Flags().VarP(&tagMap, "inject", "t", "set the value of a cue @tag field from a key=value pair")
	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", min(runtime.NumCPU(), 8), "number of concurrent build steps")

	cmd.RunE = func(c *cobra.Command, args []string) (err error) {
		path := platform
		inst, err := builder.LoadInstance(path, tagMap.Tags())
		if err != nil {
			return errors.Wrap(err)
		}

		platform, err := builder.LoadPlatform(inst)
		if err != nil {
			return errors.Wrap(err)
		}

		encoder, err := holos.NewSequentialEncoder(format, cmd.OutOrStdout())
		if err != nil {
			return errors.Wrap(err)
		}
		defer encoder.Close()

		buildPlanOpts := holos.NewBuildOpts(path)
		buildPlanOpts.Stderr = cmd.ErrOrStderr()
		buildPlanOpts.Concurrency = concurrency

		platformOpts := builder.PlatformOpts{
			Fn:          makeBuildFunc(encoder, buildPlanOpts),
			Selector:    selector,
			Concurrency: concurrency,
		}

		if err := platform.Build(c.Context(), platformOpts); err != nil {
			return errors.Wrap(err)
		}

		return nil
	}
	return cmd
}

func makeBuildFunc(encoder holos.OrderedEncoder, opts holos.BuildOpts) builder.BuildFunc {
	return func(ctx context.Context, idx int, component holos.Component) error {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err())
		default:
			inst, err := builder.LoadInstance(component.Path(), component.Tags())
			if err != nil {
				return errors.Wrap(err)
			}

			bp, err := builder.LoadBuildPlan(inst, opts)
			if err != nil {
				return errors.Wrap(err)
			}
			if err := bp.Export(idx, encoder); err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	}
}
