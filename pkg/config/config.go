package config

import (
	"flag"
	"fmt"
	"github.com/holos-run/holos/pkg/logger"
	"io"
	"log/slog"
)

// New returns a new Config
func New(stdout io.Writer, stderr io.Writer) *Config {
	return &Config{
		logConfig: logger.NewConfig(),
		stdout:    stdout,
		stderr:    stderr,
	}
}

// Config holds configuration for the whole program, used by main()
type Config struct {
	logConfig *logger.Config
	logger    *slog.Logger
	stdout    io.Writer
	stderr    io.Writer
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
	return c.logConfig.NewLogger(c.stderr)
}

func (c *Config) Stderr() io.Writer {
	return c.stderr
}
