package holos

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/holos-run/holos/internal/logger"
)

const DefaultProvisionerNamespace = "secrets"

// An Option configures a Config using [functional
// options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html).
type Option func(o *options)

type options struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	logger *slog.Logger
}

// Stdin redirects standard input to r, useful for test capture.
func Stdin(r io.Reader) Option {
	return func(o *options) { o.stdin = r }
}

// Stdout redirects standard output to w, useful for test capture.
func Stdout(w io.Writer) Option {
	return func(o *options) { o.stdout = w }
}

// Stderr redirects standard error to w, useful for test capture.
func Stderr(w io.Writer) Option {
	return func(o *options) { o.stderr = w }
}

func Logger(logger *slog.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// New returns a new top level cli Config.
func New(opts ...Option) *Config {
	cfgOptions := &options{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	for _, option := range opts {
		option(cfgOptions)
	}
	writeFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	txFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg := &Config{
		logConfig:    logger.NewConfig(),
		writeTo:      getenv("HOLOS_WRITE_TO", "deploy"),
		writeFlagSet: writeFlagSet,
		options:      cfgOptions,
		txtarFlagSet: txFlagSet,
		logger:       cfgOptions.logger,
	}
	writeFlagSet.StringVar(&cfg.writeTo, "write-to", cfg.writeTo, "write to directory")
	cfg.txtarIndex = txFlagSet.Int("index", 0, "file number to print if not 0")
	cfg.txtarQuote = txFlagSet.Bool("quote", true, "quote necessary files for use with testscript unquote")
	return cfg
}

// Config holds configuration for the whole program, used by main(). The config
// should be initialized early at a well known location in the program lifecycle
// then remain immutable.
type Config struct {
	logConfig    *logger.Config
	writeTo      string
	logger       *slog.Logger
	options      *options
	finalized    bool
	writeFlagSet *flag.FlagSet
	txtarIndex   *int
	txtarQuote   *bool
	txtarFlagSet *flag.FlagSet
}

func (c *Config) LogConfig() *logger.Config {
	return c.logConfig
}

// LogFlagSet returns the logging *flag.FlagSet for use by the command handler.
func (c *Config) LogFlagSet() *flag.FlagSet {
	return c.logConfig.FlagSet()
}

// WriteFlagSet returns a *flag.FlagSet wired to c *Config.  Useful for commands that write files.
func (c *Config) WriteFlagSet() *flag.FlagSet {
	return c.writeFlagSet
}

// TxtarFlagSet returns the *flag.FlagSet for txtar related commands.
func (c *Config) TxtarFlagSet() *flag.FlagSet {
	return c.txtarFlagSet
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

// Vet validates the config.
func (c *Config) Vet() error {
	if c == nil || c.logConfig == nil {
		return fmt.Errorf("cannot vet: not configured")
	}
	return c.logConfig.Vet()
}

// Logger returns a *slog.Logger configured by the user or the default logger if
// no logger has been configured by the user.
func (c *Config) Logger() *slog.Logger {
	if c == nil {
		return slog.Default()
	}
	if c.logger != nil {
		return c.logger
	}
	if c.logConfig == nil {
		return slog.Default()
	}
	return c.logConfig.NewLogger(c.Stderr())
}

// NewTopLevelLogger returns a *slog.Logger with a handler that filters source
// attributes. Useful as a top level error logger in main().
func (c *Config) NewTopLevelLogger() *slog.Logger {
	return c.logConfig.NewLogger(c.options.stderr)
}

// Stdin should be used instead of os.Stdin to capture input from tests.
func (c *Config) Stdin() io.Reader {
	if c == nil || c.options == nil {
		return os.Stdin
	}
	return c.options.stdin
}

// Stdout should be used instead of os.Stdout to capture output for tests.
func (c *Config) Stdout() io.Writer {
	if c == nil || c.options == nil {
		return os.Stdout
	}
	return c.options.stdout
}

// Stderr should be used instead of os.Stderr to capture output for tests.
func (c *Config) Stderr() io.Writer {
	if c == nil || c.options == nil {
		return os.Stderr
	}
	return c.options.stderr
}

// WriteTo returns the write to path configured by flags.
func (c *Config) WriteTo() string {
	if c == nil {
		return ""
	}
	return c.writeTo
}

// Printf calls fmt.Fprintf with the configured Stdout.  Errors are logged.
func (c *Config) Printf(format string, a ...any) {
	if _, err := fmt.Fprintf(c.Stdout(), format, a...); err != nil {
		c.Logger().Error("could not Fprintf", "err", err)
	}
}

// Println calls fmt.Fprintln with the configured Stdout.  Errors are logged.
func (c *Config) Println(a ...any) {
	if _, err := fmt.Fprintln(c.Stdout(), a...); err != nil {
		c.Logger().Error("could not Fprintln", "err", err)
	}
}

// Write writes to Stdout.  Errors are logged.
func (c *Config) Write(p []byte) {
	if _, err := c.Stdout().Write(p); err != nil {
		c.Logger().Error("could not write", "err", err)
	}
}

func (c *Config) TxtarQuote() bool {
	if c == nil || c.txtarIndex == nil {
		return false
	}
	return *c.txtarQuote
}

// TxtarIndex returns the
func (c *Config) TxtarIndex() int {
	if c == nil || c.txtarIndex == nil {
		return 0
	}
	return *c.txtarIndex
}

// getenv is equivalent to os.LookupEnv with a default value.
func getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
