package render

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("render")
	cmd.Args = cobra.NoArgs
	cmd.Short = "render platforms and components to manifest files"
	cmd.AddCommand(newPlatform(cfg, feature))
	cmd.AddCommand(newComponent(cfg, feature))
	return cmd
}

func newPlatform(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("platform DIRECTORY")
	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Example = "holos render platform"
	cmd.Short = "render an entire platform"

	config := client.NewConfig(cfg)
	if feature.Flag(holos.ClientFeature) {
		cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
		cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())
	}

	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", min(runtime.NumCPU(), 8), "number of components to render concurrently")
	var platform string
	cmd.Flags().StringVar(&platform, "platform", "./platform", "platform directory path")
	var selector holos.Selector
	cmd.Flags().VarP(&selector, "selector", "l", "label selector (e.g. label==string,label!=string)")
	tagMap := make(holos.TagMap)
	cmd.Flags().VarP(&tagMap, "inject", "t", "set the value of a cue @tag field from a key=value pair")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		log := logger.FromContext(ctx)
		if len(args) > 0 {
			platform = args[0]
			msg := "deprecated: %s, use the --platform flag instead"
			log.WarnContext(ctx, fmt.Sprintf(msg, platform))
		}

		inst, err := builder.LoadInstance(platform, tagMap.Tags())
		if err != nil {
			return errors.Wrap(err)
		}

		platform, err := builder.LoadPlatform(inst)
		if err != nil {
			return errors.Wrap(err)
		}

		prefixArgs := []string{
			"--log-level", cfg.LogConfig().Level(),
			"--log-format", cfg.LogConfig().Format(),
		}
		opts := builder.PlatformOpts{
			Fn:          makePlatformRenderFunc(cmd.ErrOrStderr(), prefixArgs),
			Selector:    selector,
			Concurrency: concurrency,
			InfoEnabled: true,
		}

		if err := platform.Build(ctx, opts); err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	return cmd
}

// New returns the component subcommand for the render command
func newComponent(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("component DIRECTORY")
	cmd.Args = cobra.ExactArgs(1)
	cmd.Short = "render a platform component"
	cmd.Example = "  holos render component --inject holos_cluster=aws2 ./components/monitoring/kube-prometheus-stack"
	cmd.Flags().AddGoFlagSet(cfg.WriteFlagSet())
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())

	config := client.NewConfig(cfg)
	if feature.Flag(holos.ClientFeature) {
		cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
		cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())
	}

	tagMap := make(holos.TagMap)
	cmd.Flags().VarP(&tagMap, "inject", "t", "set the value of a cue @tag field from a key=value pair")
	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", min(runtime.NumCPU(), 8), "number of concurrent build steps")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		path := args[0]

		inst, err := builder.LoadInstance(path, tagMap.Tags())
		if err != nil {
			return errors.Wrap(err)
		}

		opts := holos.NewBuildOpts(path)
		opts.Stderr = cmd.ErrOrStderr()
		opts.Concurrency = concurrency
		opts.WriteTo = cfg.WriteTo()

		bp, err := builder.LoadBuildPlan(inst, opts)
		if err != nil {
			return errors.Wrap(err)
		}

		if err := bp.Build(ctx); err != nil {
			return errors.Wrap(err)
		}

		return nil
	}
	return cmd
}

func makePlatformRenderFunc(w io.Writer, prefixArgs []string) builder.BuildFunc {
	return func(ctx context.Context, idx int, component holos.Component) error {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err())
		default:
			tags, err := component.Tags()
			if err != nil {
				return errors.Wrap(err)
			}
			args := make([]string, 0, 10+len(prefixArgs)+(len(tags)*2))
			args = append(args, prefixArgs...)
			args = append(args, "render", "component")
			for _, tag := range tags {
				args = append(args, "--inject", tag)
			}
			args = append(args, component.Path())
			if _, err := util.RunCmdW(ctx, w, "holos", args...); err != nil {
				return errors.Format("could not render component: %w", err)
			}
		}
		return nil
	}
}
