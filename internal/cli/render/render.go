package render

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"runtime"
	"strings"

	h "github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/artifact"
	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/builder/v1alpha4"
	"github.com/holos-run/holos/internal/builder/v1alpha5"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/render"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config, feature holos.Flagger) *cobra.Command {
	cmd := command.New("render")
	cmd.Args = cobra.NoArgs
	cmd.Short = "render platforms and components to manifest files"
	cmd.AddCommand(NewComponent(cfg))
	cmd.AddCommand(NewPlatform(cfg))
	return cmd
}

// New returns the component subcommand for the render command
func NewComponent(cfg *holos.Config) *cobra.Command {
	cmd := command.New("component DIRECTORY")
	cmd.Args = cobra.ExactArgs(1)
	cmd.Short = "render a platform component"
	cmd.Example = "  holos render component --inject holos_cluster=aws2 ./components/monitoring/kube-prometheus-stack"
	cmd.Flags().AddGoFlagSet(cfg.WriteFlagSet())
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	tagMap := make(tags)
	cmd.PersistentFlags().VarP(&tagMap, "inject", "t", "set the value of a cue @tag field from a key=value pair")

	var concurrency int
	flagSet.IntVar(&concurrency, "concurrency", min(runtime.NumCPU(), 8), "number of concurrent build steps")
	cmd.Flags().AddGoFlagSet(flagSet)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		build := builder.New(
			builder.Entrypoints(args),
			builder.Cluster(cfg.ClusterName()),
			builder.Tags(tagMap.Tags()),
		)
		log := logger.FromContext(ctx)

		log.DebugContext(ctx, "cue: building component instance")
		bd, err := build.Unify(ctx, config)
		if err != nil {
			return errors.Wrap(err)
		}

		typeMeta, err := bd.TypeMeta()
		if err != nil {
			return errors.Wrap(err)
		}

		if typeMeta.Kind != "BuildPlan" {
			return errors.Format("invalid kind: want: BuildPlan have: %s", typeMeta.Kind)
		}

		log.DebugContext(ctx, "discriminated "+typeMeta.APIVersion+" "+typeMeta.Kind)

		jsonBytes, err := bd.Value.MarshalJSON()
		if err != nil {
			return errors.Format("could not marshal json %s: %w", bd.Dir, err)
		}
		decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
		decoder.DisallowUnknownFields()

		art := artifact.New()

		switch version := typeMeta.APIVersion; version {
		case "v1alpha5":
			builder := v1alpha5.BuildPlan{
				WriteTo:     cfg.WriteTo(),
				Concurrency: concurrency,
				Stderr:      cmd.ErrOrStderr(),
				Path:        h.InstancePath(args[0]),
			}
			if err := decoder.Decode(&builder.BuildPlan); err != nil {
				return errors.Format("could not decode build plan %s: %w", bd.Dir, err)
			}
			return render.Component(ctx, &builder, art)
		case "v1alpha4":
			builder := v1alpha4.BuildPlan{
				WriteTo:     cfg.WriteTo(),
				Concurrency: concurrency,
				Stderr:      cmd.ErrOrStderr(),
				Path:        h.InstancePath(args[0]),
			}
			if err := decoder.Decode(&builder.BuildPlan); err != nil {
				return errors.Format("could not decode build plan %s: %w", bd.Dir, err)
			}
			return render.Component(ctx, &builder, art)
		// Legacy method.
		case "v1alpha3", "v1alpha2", "v1alpha1":
			//nolint:staticcheck
			results, err := build.Run(ctx, config)
			if err != nil {
				return errors.Wrap(err)
			}
			// TODO: Avoid accidental over-writes if two or more holos component
			// instances result in the same file path. Write files into a blank
			// temporary directory, error if a file exists, then move the directory into
			// place.
			var result Result
			for _, result = range results {
				log := logger.FromContext(ctx).With(
					"cluster", cfg.ClusterName(),
					"name", result.Name(),
				)
				if result.Continue() {
					continue
				}
				// DeployFiles from the BuildPlan
				if err := result.WriteDeployFiles(ctx, cfg.WriteTo()); err != nil {
					return errors.Wrap(err)
				}

				// API Objects
				if result.SkipWriteAccumulatedOutput() {
					log.DebugContext(ctx, "skipped writing k8s objects for "+result.Name())
				} else {
					path := result.Filename(cfg.WriteTo(), cfg.ClusterName())
					if err := result.Save(ctx, path, result.AccumulatedOutput()); err != nil {
						return errors.Wrap(err)
					}
				}

				log.InfoContext(ctx, "rendered "+result.Name(), "status", "ok", "action", "rendered")
			}

		default:
			return errors.Format("component version not supported: %s", version)
		}

		return nil
	}
	return cmd
}

func NewPlatform(cfg *holos.Config) *cobra.Command {
	cmd := command.New("platform DIRECTORY")
	cmd.Args = cobra.ExactArgs(1)
	cmd.Example = "  holos render platform ./platform"
	cmd.Short = "render an entire platform"

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", min(runtime.NumCPU(), 8), "number of components to render concurrently")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		build := builder.New(builder.Entrypoints(args))
		log := logger.FromContext(ctx)

		log.DebugContext(ctx, "cue: building platform instance")
		bd, err := build.Unify(ctx, config)
		if err != nil {
			return errors.Wrap(err)
		}

		tm, err := bd.TypeMeta()
		if err != nil {
			return errors.Wrap(err)
		}

		if tm.Kind != "Platform" {
			return errors.Format("invalid kind: want: Platform have: %s", tm.Kind)
		}

		log.DebugContext(ctx, "discriminated "+tm.APIVersion+" "+tm.Kind)

		jsonBytes, err := bd.Value.MarshalJSON()
		if err != nil {
			return errors.Format("could not marshal json %s: %w", bd.Dir, err)
		}
		decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
		decoder.DisallowUnknownFields()

		switch version := tm.APIVersion; version {
		case "v1alpha5":
			builder := v1alpha5.Platform{
				Concurrency: concurrency,
				Stderr:      cmd.ErrOrStderr(),
			}
			if err := decoder.Decode(&builder.Platform); err != nil {
				return errors.Format("could not decode platform %s: %w", bd.Dir, err)
			}
			return render.Platform(ctx, &builder)
		case "v1alpha4":
			builder := v1alpha4.Platform{
				Concurrency: concurrency,
				Stderr:      cmd.ErrOrStderr(),
			}
			if err := decoder.Decode(&builder.Platform); err != nil {
				return errors.Format("could not decode platform %s: %w", bd.Dir, err)
			}
			return render.Platform(ctx, &builder)

		// Legacy versions prior to the render.Builder interface.
		case "v1alpha3", "v1alpha2", "v1alpha1":
			platform, err := build.Platform(ctx, config)
			if err != nil {
				return errors.Wrap(err)
			}
			//nolint:staticcheck
			return render.LegacyPlatform(ctx, concurrency, platform, cmd.ErrOrStderr())

		default:
			return errors.Format("platform version not supported: %s", version)
		}
	}

	return cmd
}

// tags represents a map of key values for CUE tags for flag parsing.
type tags map[string]string

func (t tags) Tags() []string {
	parts := make([]string, 0, len(t))
	for k, v := range t {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return parts
}

func (t tags) String() string {
	return strings.Join(t.Tags(), ",")
}

func (t tags) Set(value string) error {
	for _, item := range strings.Split(value, ",") {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			return errors.Format("invalid format, must be tag=value")
		}
		t[parts[0]] = parts[1]
	}
	return nil
}

func (t tags) Type() string {
	return "strings"
}

// Deprecated: use render.Artifact instead.
type Result interface {
	Continue() bool
	Name() string
	Filename(writeTo string, cluster string) string
	KustomizationFilename(writeTo string, cluster string) string
	Save(ctx context.Context, path string, content string) error
	AccumulatedOutput() string
	SkipWriteAccumulatedOutput() bool
	WriteDeployFiles(ctx context.Context, writeTo string) error
	GetKind() string
	GetAPIVersion() string
}
