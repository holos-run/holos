package render

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/builder/v1alpha6"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const tagHelp = "set the value of a cue @tag field in the form key [ = value ]"

// typeMetaFile represents the file holos uses to discriminate the api version
// of a component BuildPlan.
const typeMetaFile = "typemeta.yaml"

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

		rp := renderParams{
			path:         args[0],
			extractYAMLs: extractYAMLs,
			tagMap:       tagMap,
			stderr:       cmd.ErrOrStderr(),
			concurrency:  concurrency,
			writeTo:      cfg.WriteTo(),
		}

		log := logger.FromContext(ctx)

		// if typemeta.yaml does not exist, render using <= v1alpha5 technique.
		typeMetaPath := filepath.Join(rp.path, typeMetaFile)
		if _, err := os.Stat(typeMetaPath); err != nil {
			log.DebugContext(ctx, fmt.Sprintf("could not load %s falling back to deprecated builder", typeMetaPath), "path", typeMetaPath, "err", err)
			return renderComponentAlpha5(ctx, rp)
		}

		// Render using the new typemeta.yaml discriminator in v1alpha6 and later.
		tm, err := loadTypeMeta(typeMetaPath)
		if err != nil {
			return errors.Format("could not discriminate component type: %w", err)
		}
		if err := renderComponent(ctx, rp, tm); err != nil {
			return errors.Format("could not render component: %w", err)
		}
		return nil
	}
	return cmd
}

func loadTypeMeta(path string) (tm holos.TypeMeta, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return tm, errors.Wrap(err)
	}
	if err = yaml.Unmarshal(data, &tm); err != nil {
		return tm, errors.Wrap(err)
	}
	return tm, nil
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

// renderParams represents parameters for the component rendering functions.
type renderParams struct {
	path         string
	extractYAMLs holos.StringSlice
	tagMap       holos.TagMap
	stderr       io.Writer
	concurrency  int
	writeTo      string
}

// renderComponent implements the behavior of holos render component for
// v1alpha6 and later component versions.  The typemeta.yaml file located in the
// component directory must be present and is used to discriminate the
// apiVersion prior to building the CUE instance.  Useful to determine which
// build tags need to be injected depending on the apiVersion of the component.
func renderComponent(ctx context.Context, rp renderParams, tm holos.TypeMeta) error {
	if tm.Kind != "BuildPlan" {
		return errors.Format("unsupported kind: %s, want BuildPlan", tm.Kind)
	}
	// so we can append version specific tags.
	tags := rp.tagMap.Tags()
	// temp directory is an important part of the build context.
	tempDir, err := os.MkdirTemp("", "holos.render")
	if err != nil {
		return errors.Format("could not make temp dir: %w", err)
	}
	defer util.Remove(ctx, tempDir)

	// Runtime configuration of the build.
	opts := holos.NewBuildOpts(rp.path)
	opts.Stderr = rp.stderr
	opts.Concurrency = rp.concurrency
	opts.WriteTo = rp.writeTo
	opts.BuildContext.TempDir = tempDir

	log := logger.FromContext(ctx)
	log.DebugContext(ctx, fmt.Sprintf("rendering %s kind %s version %s", rp.path, tm.Kind, tm.APIVersion), "kind", tm.Kind, "apiVersion", tm.APIVersion, "path", rp.path)

	// generic build plan wrapper for all versions.
	var bp builder.BuildPlan
	switch version := tm.APIVersion; version {
	case "v1alpha6":
		// Prepare runtime build context for injection as a cue tag.
		bc := v1alpha6.NewBuildContext(opts.BuildContext)
		buildContextTags, err := bc.Tags()
		if err != nil {
			return errors.Format("could not get build context tag: %w", err)
		}
		tags = append(tags, buildContextTags...)
		// the version specific build plan itself embedded into the wrapper.
		bp = builder.BuildPlan{BuildPlan: &v1alpha6.BuildPlan{Opts: opts}}
	default:
		return errors.Format("unsupported version: %s, must be at least v1alpha6 when typemeta.yaml is present", version)
	}

	// Load the CUE instance to export the BuildPlan with version-specific tags.
	inst, err := builder.LoadInstance(rp.path, rp.extractYAMLs, tags)
	if err != nil {
		return errors.Format("could not load cue instance: %w", err)
	}

	// Get the holos: field value from cue.
	v, err := inst.HolosValue()
	if err != nil {
		return errors.Wrap(err)
	}

	// Load the BuildPlan from the cue value.
	if err := bp.Load(v); err != nil {
		return errors.Wrap(err)
	}

	// Execute the build.
	if err := bp.Build(ctx); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// renderComponentAlpha5 implements the behavior of holos render component for
// v1alpha5 and earlier.  This method loads the CUE Instance to discriminate the
// apiVersion, which is too late to pass tags properly.
//
// Deprecated: use renderComponent() instead
func renderComponentAlpha5(ctx context.Context, rp renderParams) error {
	// Manage a temp directory for the build artifacts.  The concrete value is
	// needed prior to exporting the BuildPlan from the CUE instance.
	tempDir, err := os.MkdirTemp("", "holos.render")
	if err != nil {
		return errors.Format("could not make temp dir: %w", err)
	}
	defer util.Remove(ctx, tempDir)

	// Load the CUE instance to export the BuildPlan.
	inst, err := builder.LoadInstance(rp.path, rp.extractYAMLs, rp.tagMap.Tags())
	if err != nil {
		return errors.Format("could not load cue instance: %w", err)
	}

	// Runtime configuration of the build.
	opts := holos.NewBuildOpts(rp.path)
	opts.Stderr = rp.stderr
	opts.Concurrency = rp.concurrency
	opts.WriteTo = rp.writeTo

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
