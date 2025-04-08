package render

import (
	"context"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/component"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/platform"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
)

func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("render")
	cmd.Args = cobra.NoArgs
	cmd.Short = "render platforms and components to manifest files"
	cmd.AddCommand(NewRenderPlatformCommand(cfg, platform.NewConfig()))
	cmd.AddCommand(component.NewCommand(component.NewConfig()))
	return cmd
}

func NewRenderPlatformCommand(cfg *holos.Config, pcfg *platform.Config) (cmd *cobra.Command) {
	rp := &renderPlatform{cfg: cfg, pcfg: pcfg}
	cmd = platform.NewCommand(pcfg, rp.Run)
	cmd.Short = "render an entire platform"
	cmd.Flags().AddFlagSet(pcfg.FlagSet())
	return cmd
}

// renderPlatform implements the holos render platform command.
type renderPlatform struct {
	cfg  *holos.Config
	pcfg *platform.Config
}

// Run executes the holos render component command concurrently for each
// platform component.  The overall approach is to marshal the component into
// cue tags, pass the log level and format, then execute the command as a sub
// process.
//
// Note that the marshalling of component fields through the argument vector has
// been quite awkward to maintain.  Consider refactoring to an approach of
// marshaling the entire component structure via stdin to a sub process that
// returns a build plan on stdout.  The purpose of using a sub process is to
// execute cue concurrently.  Cue is not safe for concurrent use within the same
// process.
func (r *renderPlatform) Run(ctx context.Context, p *platform.Platform) error {
	prefixArgs := []string{
		"--log-level", r.cfg.LogConfig().Level(),
		"--log-format", r.cfg.LogConfig().Format(),
	}

	opts := platform.BuildOpts{
		PerComponentFunc: func(ctx context.Context, i int, c holos.Component) error {
			select {
			case <-ctx.Done():
				return errors.Wrap(ctx.Err())
			default:
				args := make([]string, 0, 100)
				args = append(args, prefixArgs...)
				args = append(args, "render", "component")
				// holos render platform --inject tags
				for _, tag := range r.pcfg.TagMap.Tags() {
					args = append(args, "--inject", tag)
				}
				// component tags (name, labels, annotations)
				if tags, err := c.Tags(); err != nil {
					return errors.Wrap(err)
				} else {
					for _, tag := range tags {
						args = append(args, "--inject", tag)
					}
				}
				// component path
				args = append(args, c.Path())

				// Get current executable path.
				holosPath, err := util.Executable()
				if err != nil {
					return errors.Wrap(err)
				}

				// Run holos render component ...
				if _, err := util.RunCmdA(ctx, r.pcfg.Stderr, holosPath, args...); err != nil {
					return errors.Format("could not render component: %w", err)
				}
			}
			return nil
		},
		InfoEnabled: true,
	}

	return errors.Wrap(p.Build(ctx, opts))
}
