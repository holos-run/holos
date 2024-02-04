// Package logger provides logging configuration and helpers to pass a logger instance through the context.
package logger

import (
	"context"
	"flag"
	"fmt"
	"github.com/holos-run/holos/pkg/version"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const ErrKey = "err"

var validLogLevels = []string{"debug", "info", "warn", "error"}
var validLogFormats = []string{"text", "json"}

// stringSlice is a comma separated list of string values
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join((*s)[:], ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}

// key is an unexported type for keys defined in this package to prevent
// collisions with keys defined in other packages.
type key int

// https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/context/context.go;l=140-158
// loggerKey is the key for *slog.Logs values in Contexts. It us unexported;
// clients use NewContext and FromContext instead of this key directly.
var loggerKey key

// NewContext returns a new Context that carries value logger. Use FromContext
// to retrieve the value.
func NewContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the *slog.Logs previously stored in ctx by NewContext.
// slog.Default() is returned otherwise.
func FromContext(ctx context.Context) *slog.Logger {
	// https://go.dev/ref/spec#Type_assertions
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// Config specifies user configurable flag values to create a NewLogger
type Config struct {
	level     string
	format    string
	dropAttrs stringSlice
	flagSet   *flag.FlagSet
}

// GetLogLevel returns a slog.Level configured by the user
//
// A non-zero length DEBUG env var takes precedence over config fields.
func (c *Config) GetLogLevel() slog.Level {
	if os.Getenv("DEBUG") != "" {
		return slog.LevelDebug
	}
	switch strings.ToLower(c.level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (c *Config) ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if slices.Contains(c.dropAttrs, a.Key) {
		return slog.Attr{}
	}
	// Check if err
	if a.Key == ErrKey {
		if err, ok := a.Value.Any().(error); ok {
			return tint.Err(err)
		}
		if err, ok := a.Value.Any().(string); ok {
			return tint.Err(fmt.Errorf(err))
		}
	} else if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
	}
	return a
}

// NewTopLevelLogger returns a *slog.Logs configured by c *Config which writes
// to w without source information. Useful as a top level logger where the
// source is know and the error is wrapped with a location attribute.
func (c *Config) NewTopLevelLogger(w io.Writer) *slog.Logger {
	level := c.GetLogLevel()
	var handler slog.Handler
	if c.format == "text" {
		noColor := true
		if file, ok := w.(*os.File); ok {
			noColor = !isatty.IsTerminal(file.Fd())
		}
		handler = tint.NewHandler(w, &tint.Options{
			Level:       level,
			TimeFormat:  time.Kitchen,
			AddSource:   false,
			ReplaceAttr: c.ReplaceAttr,
			NoColor:     noColor,
		})
	} else {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:       level,
			AddSource:   false,
			ReplaceAttr: c.ReplaceAttr,
		})
	}

	return slog.New(handler).With("version", version.Version)
}

// NewLogger returns a *slog.Logs configured by c *Config which writes to w
func (c *Config) NewLogger(w io.Writer) *slog.Logger {
	level := c.GetLogLevel()
	var handler slog.Handler
	if c.format == "text" {
		noColor := true
		if file, ok := w.(*os.File); ok {
			noColor = !isatty.IsTerminal(file.Fd())
		}
		handler = tint.NewHandler(w, &tint.Options{
			Level:       level,
			TimeFormat:  time.Kitchen,
			AddSource:   true,
			ReplaceAttr: c.ReplaceAttr,
			NoColor:     noColor,
		})
	} else {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:       level,
			AddSource:   true,
			ReplaceAttr: c.ReplaceAttr,
		})
	}

	return slog.New(handler).With("version", version.Version)
}

// NewConfig returns a new logging Config struct
func NewConfig() *Config {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	c := &Config{flagSet: f}
	f.StringVar(&c.level, "log-level", "info", fmt.Sprintf("Log Level (%s)", strings.Join(validLogLevels, "|")))
	f.StringVar(&c.format, "log-format", "text", fmt.Sprintf("Log format (%s)", strings.Join(validLogFormats, "|")))
	f.Var(&c.dropAttrs, "log-drop", "Log attributes to drop, e.g. \"user-agent,version\"")
	return c
}

// FlagSet returns the go flag set to configure logging
func (c *Config) FlagSet() *flag.FlagSet {
	return c.flagSet
}

// Vet validates the config values
func (c *Config) Vet() error {
	if err := c.vetLevel(); err != nil {
		return err
	}
	if err := c.vetFormat(); err != nil {
		return err
	}
	return nil
}

func (c *Config) vetLevel() error {
	for _, validLevel := range validLogLevels {
		if c.level == validLevel {
			return nil
		}
	}
	err := fmt.Errorf("invalid log level: %s is not one of %s", c.level, strings.Join(validLogLevels, ", "))
	return wrapper.Wrap(err)
}

func (c *Config) vetFormat() error {
	for _, validFormat := range validLogFormats {
		if c.format == validFormat {
			return nil
		}
	}
	err := fmt.Errorf("invalid log format: %s is not one of %s", c.format, strings.Join(validLogFormats, ", "))
	return wrapper.Wrap(err)
}
