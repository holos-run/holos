package platform

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/cue"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/platform/v1alpha5"
	"github.com/holos-run/holos/internal/platform/v1alpha6"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

func New(cfg Config, root, leaf string) *Platform {
	return &Platform{
		cfg:  cfg,
		root: root,
		leaf: leaf,
	}
}

func NewCommand(cfg Config, run func(context.Context, *Platform) error) *cobra.Command {
	cmd := command.New("platform")
	cmd.Short = "process a platform resource"
	cmd.Args = cobra.MinimumNArgs(0)
	cmd.Flags().AddFlagSet(cfg.flagSet())
	cmd.SetOut(cfg.Stdout)
	cmd.SetErr(cfg.Stderr)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		wd, err := os.Getwd()
		if err != nil {
			return errors.Format("could not get current directory: %w", err)
		}
		if len(args) == 0 {
			args = append(args, "platform")
		}
		for _, leaf := range args {
			if filepath.IsAbs(leaf) {
				wd, leaf, err = util.FindRootLeaf(leaf)
				if err != nil {
					return errors.Wrap(err)
				}
			}
			p := New(cfg, wd, leaf)
			if err := p.Load(ctx); err != nil {
				return errors.Format("could not load %s: %w", leaf, err)
			}
			if err := run(ctx, p); err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	}
	return cmd
}

func NewConfig() Config {
	cfg := Config{
		Concurrency: runtime.NumCPU(),
		TagMap:      make(holos.TagMap),
		WriteTo:     os.Getenv(holos.WriteToEnvVar),
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
	}
	if cfg.WriteTo == "" {
		cfg.WriteTo = holos.WriteToDefault
	}
	return cfg
}

type Config struct {
	// TagMap represents cue tags to inject.
	TagMap holos.TagMap
	// Concurrency represents the number of subcommands to execute concurrently.
	Concurrency int
	// WriteTo represents the output base directory for rendered artifacts.
	WriteTo string
	// Stdout represents the standard output pipe.
	Stdout io.Writer
	// Stderr represents the standard error pipe.  Used to copy stderr output from
	// subcommands.
	Stderr io.Writer
}

func (c *Config) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&c.WriteTo, "write-to", c.WriteTo, fmt.Sprintf("write to directory (%s)", holos.WriteToEnvVar))
	fs.VarP(c.TagMap, "inject", "t", holos.TagMapHelp)
	fs.IntVar(&c.Concurrency, "concurrency", c.Concurrency, "number of concurrent build steps")
	return fs
}

type Platform struct {
	holos.Platform
	cfg  Config
	root string
	leaf string
}

// Root returns the platform root directory.
func (p *Platform) Root() string {
	return p.root
}

// Load discriminates the api version then loads the platform configuration by
// building a cue instance.
func (p *Platform) Load(ctx context.Context) error {
	tags := p.cfg.TagMap.Tags()

	tm, err := cue.TypeMeta(p.root, p.leaf)
	if err != nil {
		return errors.Wrap(err)
	}

	switch tm.APIVersion {
	case "v1alpha6":
		p.Platform = &v1alpha6.Platform{}
	default:
		p.Platform = &v1alpha5.Platform{}
	}

	inst, err := cue.BuildInstance(p.root, p.leaf, tags)
	if err != nil {
		return errors.Format("could not build cue instance: %w", err)
	}
	val, err := inst.HolosValue()
	if err != nil {
		return errors.Format("could not get holos field value: %w", err)
	}
	if err := p.Platform.Load(val); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// BuildOpts represents build options when processing the components in a
// platform.
type BuildOpts struct {
	PerComponentFunc   func(context.Context, int, holos.Component) error
	ComponentSelectors holos.Selectors
	InfoEnabled        bool
}

// Build calls [opts.PerComponentFunc] for each platform component.
func (p *Platform) Build(ctx context.Context, opts BuildOpts) error {
	limit := max(1, p.cfg.Concurrency)
	parentStart := time.Now()
	components := p.Select(opts.ComponentSelectors...)
	total := len(components)

	g, ctx := errgroup.WithContext(ctx)
	// Limit the number of concurrent goroutines due to CUE memory usage concerns
	// while rendering components.  One more for the producer.
	g.SetLimit(limit + 1)
	// Spawn a producer because g.Go() blocks when the group limit is reached.
	g.Go(func() error {
		for idx := range components {
			select {
			case <-ctx.Done():
				return errors.Wrap(ctx.Err())
			default:
				// Capture idx to avoid issues with closure.  Fixed in Go 1.22.
				idx := idx
				component := components[idx]
				// Worker go routine. Blocks if limit has been reached.
				g.Go(func() error {
					select {
					case <-ctx.Done():
						return errors.Wrap(ctx.Err())
					default:
						start := time.Now()
						log := logger.FromContext(ctx).With("num", idx+1, "total", total)
						if err := opts.PerComponentFunc(ctx, idx, component); err != nil {
							return errors.Wrap(err)
						}
						duration := time.Since(start)
						msg := fmt.Sprintf("rendered %s in %s", component.Describe(), duration)
						if opts.InfoEnabled {
							log.InfoContext(ctx, msg, "duration", duration)
						} else {
							log.DebugContext(ctx, msg, "duration", duration)
						}
						return nil
					}
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Since(parentStart)
	msg := fmt.Sprintf("rendered platform in %s", duration)
	if opts.InfoEnabled {
		logger.FromContext(ctx).InfoContext(ctx, msg, "duration", duration)
	} else {
		logger.FromContext(ctx).DebugContext(ctx, msg, "duration", duration)
	}
	return nil
}
