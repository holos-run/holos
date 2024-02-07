package config

import (
	"flag"
	"fmt"
	"github.com/holos-run/holos/pkg/logger"
	"io"
	"log/slog"
	"os"
)

// An Option configures a Config
type Option func(o *options)

type options struct {
	stdout io.Writer
	stderr io.Writer
}

// Stdout redirects standard output to w
func Stdout(w io.Writer) Option {
	return func(o *options) { o.stdout = w }
}

// Stderr redirects standard error to w
func Stderr(w io.Writer) Option {
	return func(o *options) { o.stderr = w }
}

// New returns a new top level cli Config.
func New(opts ...Option) *Config {
	o := &options{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	for _, f := range opts {
		f(o)
	}
	return &Config{
		logConfig: logger.NewConfig(),
		options:   o,
	}
}

// Config holds configuration for the whole program, used by main()
type Config struct {
	logConfig *logger.Config
	logger    *slog.Logger
	options   *options
	finalized bool
}

// LogFlagSet returns the logging *flag.FlagSet
func (c *Config) LogFlagSet() *flag.FlagSet {
	return c.logConfig.FlagSet()
}

// Finalize validates the config and finalizes the startup lifecycle based on user configuration.
func (c *Config) Finalize() error {
	if c.finalized {
		return fmt.Errorf("could not finalize: already finalized")
	}
	if err := c.Vet(); err != nil {
		return err
	}
	l := c.Logger()
	c.logger = l
	l.Debug("config lifecycle", "state", "finalized")
	c.finalized = true
	return nil
}

// Vet validates the config
func (c *Config) Vet() error {
	return c.logConfig.Vet()
}

// Logger returns a *slog.Logger configured by the user
func (c *Config) Logger() *slog.Logger {
	if c.logger != nil {
		return c.logger
	}
	return c.logConfig.NewLogger(c.options.stderr)
}

// NewTopLevelLogger returns a *slog.Logger with a handler that filters source
// attributes. Useful as a top level error logger in main().
func (c *Config) NewTopLevelLogger() *slog.Logger {
	return c.logConfig.NewTopLevelLogger(c.options.stderr)
}

func (c *Config) Stderr() io.Writer {
	return c.options.stderr
}
