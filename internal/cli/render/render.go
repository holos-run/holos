package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
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

const buildContextFile = "components/buildcontext.json"

const tagHelp = "set the value of a cue @tag field in the form key [ = value ]"

func New(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("render")
	cmd.Args = cobra.NoArgs
	cmd.Short = "render platforms and components to manifest files"
	cmd.AddCommand(newPlatform(cfg, feature))
	cmd.AddCommand(newComponent(cfg, feature))
	return cmd
}

func newPlatform(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("platform")
	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Example = "holos render platform"
	cmd.Short = "render an entire platform"

	config := client.NewConfig(cfg)
	if feature.Flag(holos.ClientFeature) {
		cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
		cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())
	}

	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", runtime.NumCPU(), "number of components to render concurrently")
	var platform string
	cmd.Flags().StringVar(&platform, "platform", "./platform", "platform directory path")
	var extractYAMLs holos.StringSlice
	cmd.Flags().Var(&extractYAMLs, "extract-yaml", "data file paths to extract and unify with the platform config")
	var selectors holos.Selectors
	cmd.Flags().VarP(&selectors, "selector", "l", "label selector (e.g. label==string,label!=string)")
	tagMap := make(holos.TagMap)
	cmd.Flags().VarP(&tagMap, "inject", "t", tagHelp)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		log := logger.FromContext(ctx)
		if len(args) > 0 {
			platform = args[0]
			msg := "deprecated: %s, use the --platform flag instead"
			log.WarnContext(ctx, fmt.Sprintf(msg, platform))
		}

		inst, err := builder.LoadInstance(platform, extractYAMLs, tagMap.Tags())
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
			Fn:          makeComponentRenderFunc(cmd.ErrOrStderr(), prefixArgs, tagMap.Tags()),
			Selectors:   selectors,
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
	cmd.Flags().VarP(&tagMap, "inject", "t", tagHelp)
	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", runtime.NumCPU(), "number of concurrent build steps")
	var extractYAMLs holos.StringSlice
	cmd.Flags().Var(&extractYAMLs, "extract-yaml", "data file paths to extract and unify with the platform config")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		path := args[0]

		// Manage a temp directory for the build artifacts.  The concrete value is
		// needed prior to exporting the BuildPlan from the CUE instance.
		tempDir, err := os.MkdirTemp("", "holos.render")
		if err != nil {
			return errors.Format("could not make temp dir: %w", err)
		}
		defer util.Remove(ctx, tempDir)

		// Write the build context to the filesystem for cue embed.
		if err := saveBuildContext(ctx, tempDir); err != nil {
			return errors.Format("could not save build context: %w", err)
		}

		// Load the CUE instance to export the BuildPlan.
		inst, err := builder.LoadInstance(path, extractYAMLs, tagMap.Tags())
		if err != nil {
			return errors.Format("could not load cue instance: %w", err)
		}

		// Runtime configuration of the build.
		opts := holos.NewBuildOpts(path)
		opts.Stderr = cmd.ErrOrStderr()
		opts.Concurrency = concurrency
		opts.WriteTo = cfg.WriteTo()

		// Export the BuildPlan from the CUE instance.
		bp, err := builder.LoadBuildPlan(inst, opts)
		if err != nil {
			return errors.Wrap(err)
		}

		// Execute the build.
		if err := bp.Build(ctx); err != nil {
			return errors.Wrap(err)
		}

		return nil
	}
	return cmd
}

func makeComponentRenderFunc(w io.Writer, prefixArgs, cliTags []string) func(context.Context, int, holos.Component) error {
	return func(ctx context.Context, idx int, component holos.Component) error {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err())
		default:
			tags, err := component.Tags()
			if err != nil {
				return errors.Wrap(err)
			}
			filepaths, err := component.ExtractYAML()
			if err != nil {
				return errors.Wrap(err)
			}
			args := make([]string, 0, 10+len(prefixArgs)+(len(tags)*2+len(filepaths)*2))
			args = append(args, prefixArgs...)
			args = append(args, "render", "component")
			for _, tag := range cliTags {
				args = append(args, "--inject", tag)
			}
			for _, tag := range tags {
				args = append(args, "--inject", tag)
			}
			for _, path := range filepaths {
				args = append(args, "--extract-yaml", path)
			}
			args = append(args, component.Path())
			if _, err := util.RunCmdA(ctx, w, "holos", args...); err != nil {
				return errors.Format("could not render component: %w", err)
			}
		}
		return nil
	}
}

func saveBuildContext(ctx context.Context, tempDir string) error {
	buildContext := holos.BuildContext{TempDir: tempDir}

	if err := os.MkdirAll(path.Dir(buildContextFile), 0777); err != nil {
		return errors.Format("could not make components directory: %w", err)
	}

	data, err := json.MarshalIndent(buildContext, "", "  ")
	if err != nil {
		return errors.Format("could not marshal build context: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(buildContextFile, data, 0666); err != nil {
		return errors.Format("could not write build context file: %w", err)
	}

	log := logger.FromContext(ctx)
	log.DebugContext(ctx, fmt.Sprintf("wrote build context to %s", buildContextFile), "path", buildContextFile)
	return nil
}
