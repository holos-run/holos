package app

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/holos-run/holos/pkg/errors"
	"github.com/holos-run/holos/pkg/version"
	"gopkg.in/yaml.v3"
)

type Option func(*Config)

func WithIssuer(issuer string) Option {
	return func(c *Config) {
		c.OIDCIssuer = issuer
	}
}

func NewConfig(options ...Option) *Config {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	// default options
	c := &Config{
		OIDCIssuer: "https://login.ois.run",
		flagSet:    f,
	}
	for _, o := range options {
		o(c)
	}
	f.BoolVar(&c.PrintVersion, "version", false, "print version and exit")
	f.BoolVar(&c.PrintVersionYAML, "version-detail", false, "print detailed version info and exit")
	f.StringVar(&c.LogLevel, "log-level", "info", fmt.Sprintf("Log Level (%s)", strings.Join(validLogLevels, "|")))
	f.StringVar(&c.LogFormat, "log-format", "json", fmt.Sprintf("Log format (%s)", strings.Join(validLogFormats, "|")))
	f.StringVar(&c.OIDCIssuer, "oidc-issuer", c.OIDCIssuer, "OIDC Issuer URL")
	f.Var(&c.OIDCAudiences, "oidc-audience", "oidc audience to allow, e.g. \"holos-cli,https://sso.holos.run\"")
	f.BoolVar(&c.ListenAndServe, "serve", true, "Listen and serve requests.")
	f.Var(&c.dropAttrs, "log-drop", "Log attributes to drop, e.g. \"user-agent,version\"")
	f.StringVar(&c.dbURIFile, "db-uri-file", "", "File path containing the database uri")
	f.IntVar(&c.listenPort, "listen-port", 3000, "Primary service listening port")
	f.IntVar(&c.metricsPort, "metrics-port", 9090, "Prometheus metrics listen port")
	return c
}

// Config specifies user configurable values
type Config struct {
	LogLevel         string      // --log-level
	LogFormat        string      // --log-format
	PrintVersion     bool        // --version
	PrintVersionYAML bool        // --version-detail
	OIDCIssuer       string      // --oidc-issuer
	OIDCAudiences    stringSlice // --oidc-audience
	ListenAndServe   bool        // --no-serve
	dropAttrs        stringSlice // --log-drop
	dbURIFile        string      // --db-uri-file
	listenPort       int         // --listen-port
	metricsPort      int         // --metrics-port
	databaseURI      string
	flagSet          *flag.FlagSet
}

// DatabaseURI represents the database connection uri.
func (c *Config) DatabaseURI() string {
	return c.databaseURI
}

// ListenPort returns the port of the main server.
func (c *Config) ListenPort() int {
	return c.listenPort
}

// MetricsPort returns the port of the prometheus /metrics scrape endpoint configured by a flag.
func (c *Config) MetricsPort() int {
	return c.metricsPort
}

func (c *Config) FlagSet() *flag.FlagSet {
	return c.flagSet
}

func (c *Config) ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if slices.Contains(c.dropAttrs, a.Key) {
		return slog.Attr{}
	}
	return a
}

// Validate validates the config
func (c *Config) Validate(out io.Writer) error {
	if c.PrintVersion {
		return PrintVersion(out)
	}
	if c.PrintVersionYAML {
		return PrintVersionYAML(out)
	}
	if err := c.ValidateLogLevel(); err != nil {
		return errors.Wrap(err)
	}
	if err := c.ValidateLogFormat(); err != nil {
		return errors.Wrap(err)
	}
	if !strings.HasPrefix(c.OIDCIssuer, "https://") {
		return errors.Wrap(errors.New("oidc issuer must start with https://"))
	}
	if c.dbURIFile == "" {
		c.databaseURI = os.Getenv("DATABASE_URL")
	} else {
		dat, err := os.ReadFile(c.dbURIFile)
		if err != nil {
			return errors.Wrap(fmt.Errorf("could not read db uri file: %w", err))
		}
		c.databaseURI = strings.TrimSpace(string(dat))
	}
	if c.databaseURI == "" {
		slog.Warn("no database url, set DATABASE_URL env var or --db-uri-file flag")
	}
	return nil
}

func (c *Config) ValidateLogLevel() error {
	for _, validLevel := range validLogLevels {
		if c.LogLevel == validLevel {
			return nil
		}
	}
	err := fmt.Errorf("invalid log level: %s is not one of %s", c.LogLevel, strings.Join(validLogLevels, ", "))
	return errors.Wrap(err)
}

func (c *Config) ValidateLogFormat() error {
	for _, validFormat := range validLogFormats {
		if c.LogFormat == validFormat {
			return nil
		}
	}
	err := fmt.Errorf("invalid log format: %s is not one of %s", c.LogFormat, strings.Join(validLogFormats, ", "))
	return errors.Wrap(err)
}

// GetLogLevel returns a slog.Level configured by the user
//
// A non-zero length DEBUG env var takes precedence over config fields.
func (c *Config) GetLogLevel() slog.Level {
	if os.Getenv("DEBUG") != "" {
		return slog.LevelDebug
	}
	switch strings.ToLower(c.LogLevel) {
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

// PrintVersion prints the short version string only with no leading "v" to out.
func PrintVersion(out io.Writer) error {
	info := version.NewVersionInfo()
	if _, err := out.Write([]byte(info.Version + "\n")); err != nil {
		return errors.Wrap(err)
	}
	return &FastExitError{}
}

// PrintVersionYAML prints version info as yaml to out.
func PrintVersionYAML(out io.Writer) error {
	info := version.NewVersionInfo()
	data, err := yaml.Marshal(info)
	if err != nil {
		return errors.Wrap(err)
	}
	if _, err = out.Write(data); err != nil {
		return errors.Wrap(err)
	}
	return &FastExitError{}
}

type FastExitError struct{}

func (f *FastExitError) Error() string {
	return "fast exit for --version flag, this error should be handled by main()"
}
