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
	cfgOptions := &options{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	for _, option := range opts {
		option(cfgOptions)
	}
	writeFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	clusterFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg := &Config{
		logConfig:      logger.NewConfig(),
		writeTo:        getenv("HOLOS_WRITE_TO", "deploy"),
		clusterName:    getenv("HOLOS_CLUSTER_NAME", ""),
		writeFlagSet:   writeFlagSet,
		clusterFlagSet: clusterFlagSet,
		options:        cfgOptions,
	}
	writeFlagSet.StringVar(&cfg.writeTo, "write-to", cfg.writeTo, "write to directory")
	clusterFlagSet.StringVar(&cfg.clusterName, "cluster-name", cfg.clusterName, "cluster name")
	return cfg
}

// Config holds configuration for the whole program, used by main()
type Config struct {
	logConfig      *logger.Config
	writeTo        string
	clusterName    string
	logger         *slog.Logger
	options        *options
	finalized      bool
	writeFlagSet   *flag.FlagSet
	clusterFlagSet *flag.FlagSet
}

// LogFlagSet returns the logging *flag.FlagSet
func (c *Config) LogFlagSet() *flag.FlagSet {
	return c.logConfig.FlagSet()
}

// WriteFlagSet returns a *flag.FlagSet wired to c *Config.  Useful for commands that write files.
func (c *Config) WriteFlagSet() *flag.FlagSet {
	return c.writeFlagSet
}

// ClusterFlagSet returns a *flag.FlagSet wired to c *Config.  Useful for commands scoped to one cluster.
func (c *Config) ClusterFlagSet() *flag.FlagSet {
	return c.clusterFlagSet
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
	l.Debug("finalized config from flags", "state", "finalized")
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

// Stderr should be used instead of os.Stderr to capture output for tests.
func (c *Config) Stderr() io.Writer {
	return c.options.stderr
}

// Stdout should be used instead of os.Stdout to capture output for tests.
func (c *Config) Stdout() io.Writer {
	return c.options.stdout
}

// WriteTo returns the write to path configured by flags.
func (c *Config) WriteTo() string {
	return c.writeTo
}

// ClusterName returns the configured cluster name
func (c *Config) ClusterName() string {
	return c.clusterName
}

// getenv is equivalent to os.Getenv() with a default value
func getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
