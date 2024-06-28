package render

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/render"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("render")
	cmd.Args = cobra.NoArgs
	cmd.Short = "render platform configuration"
	cmd.AddCommand(NewComponent(cfg))
	cmd.AddCommand(NewPlatform(cfg))
	return cmd
}

// New returns the component subcommand for the render command
func NewComponent(cfg *holos.Config) *cobra.Command {
	cmd := command.New("component [directory...]")
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Short = "write kubernetes api objects to the filesystem"
	cmd.Flags().AddGoFlagSet(cfg.WriteFlagSet())
	cmd.Flags().AddGoFlagSet(cfg.ClusterFlagSet())

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	var printInstances bool
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.BoolVar(&printInstances, "print-instances", false, "expand /... paths for xargs")
	cmd.Flags().AddGoFlagSet(flagSet)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		build := builder.New(builder.Entrypoints(args), builder.Cluster(cfg.ClusterName()))

		if printInstances {
			instances, err := build.Instances(ctx, config)
			if err != nil {
				return errors.Wrap(err)
			}
			for _, instance := range instances {
				fmt.Fprintln(cmd.OutOrStdout(), instance.Dir)
			}
			return nil
		}

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
			// Build plans don't have anything but DeployFiles to write.
			if result.GetKind() == "BuildPlan" {
				continue
			}

			// API Objects
			path := result.Filename(cfg.WriteTo(), cfg.ClusterName())
			if err := result.Save(ctx, path, result.AccumulatedOutput()); err != nil {
				return errors.Wrap(err)
			}

			log.InfoContext(ctx, "rendered "+result.Name(), "status", "ok", "action", "rendered")
		}
		return nil
	}
	return cmd
}

func NewPlatform(cfg *holos.Config) *cobra.Command {
	cmd := command.New("platform [directory]")
	cmd.Args = cobra.ExactArgs(1)
	cmd.Short = "render all platform components"

	config := client.NewConfig(cfg)
	cmd.PersistentFlags().AddGoFlagSet(config.ClientFlagSet())
	cmd.PersistentFlags().AddGoFlagSet(config.TokenFlagSet())

	var concurrency int
	cmd.Flags().IntVar(&concurrency, "concurrency", min(runtime.NumCPU(), 8), "Number of concurrent components to render")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		build := builder.New(builder.Entrypoints(args))

		platform, err := build.Platform(ctx, config)
		if err != nil {
			return errors.Wrap(err)
		}

		return render.Platform(ctx, concurrency, platform, cmd.ErrOrStderr())
	}

	return cmd
}

type Result interface {
	Continue() bool
	Name() string
	Filename(writeTo string, cluster string) string
	KustomizationFilename(writeTo string, cluster string) string
	Save(ctx context.Context, path string, content string) error
	AccumulatedOutput() string
	WriteDeployFiles(ctx context.Context, writeTo string) error
	GetKind() string
	GetAPIVersion() string
}
