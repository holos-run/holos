package component

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Config struct {
	// TagMap represents cue tags to inject.
	TagMap holos.TagMap
	// Concurrency represents the number of subcommands to execute concurrently.
	Concurrency int
	// WriteTo represents the output base directory for rendered artifacts.
	WriteTo string
	// Stderr represents the standard error output pipe.  Used to copy stderr
	// output from subcommands.
	Stderr io.Writer
}

func (c *Config) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&c.WriteTo, "write-to", c.WriteTo, fmt.Sprintf("write to directory (%s)", holos.WriteToEnvVar))
	fs.VarP(c.TagMap, "inject", "t", holos.TagMapHelp)
	fs.IntVar(&c.Concurrency, "concurrency", c.Concurrency, "number of concurrent build steps")
	return fs
}

func NewConfig() Config {
	cfg := Config{
		Concurrency: runtime.NumCPU(),
		TagMap:      make(holos.TagMap),
		Stderr:      os.Stderr,
		WriteTo:     os.Getenv(holos.WriteToEnvVar),
	}
	if cfg.WriteTo == "" {
		cfg.WriteTo = holos.WriteToDefault
	}
	return cfg
}

// NewCommand returns the subcommand that renders individual components.
func NewCommand(cfg Config) *cobra.Command {
	cmd := command.New("component DIRECTORY")
	cmd.Args = cobra.ExactArgs(1)
	cmd.Short = "render a platform component"
	cmd.Flags().AddFlagSet(cfg.flagSet())
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Root().Context()
		// TODO(jjm): Handle fully qualified paths for tests where cwd != tempdir
		root, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err)
		}
		component := New(root, args[0], cfg)
		return component.Render(ctx)
	}
	return cmd
}
