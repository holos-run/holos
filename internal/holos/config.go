package holos

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const DefaultProvisionerNamespace = "secrets"

// An Option configures a Config using [functional
// options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html).
type Option func(o *options)

type options struct {
	stdin                io.Reader
	stdout               io.Writer
	stderr               io.Writer
	provisionerClientset kubernetes.Interface
	clientset            kubernetes.Interface
	logger               *slog.Logger
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

// ProvisionerClientset sets the kubernetes Clientset, useful for test fake.
func ProvisionerClientset(clientset kubernetes.Interface) Option {
	return func(o *options) { o.provisionerClientset = clientset }
}

// ClusterClientset sets the kubernetes Clientset, useful for test fake.
func ClusterClientset(clientset *kubernetes.Clientset) Option {
	return func(o *options) { o.clientset = clientset }
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
	clusterFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	kvFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	txFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg := &Config{
		logConfig:            logger.NewConfig(),
		writeTo:              getenv("HOLOS_WRITE_TO", "deploy"),
		clusterName:          getenv("HOLOS_CLUSTER_NAME", ""),
		writeFlagSet:         writeFlagSet,
		clusterFlagSet:       clusterFlagSet,
		options:              cfgOptions,
		kvFlagSet:            kvFlagSet,
		txtarFlagSet:         txFlagSet,
		provisionerClientset: cfgOptions.provisionerClientset,
		logger:               cfgOptions.logger,
		ServerConfig:         &ServerConfig{},
	}
	writeFlagSet.StringVar(&cfg.writeTo, "write-to", cfg.writeTo, "write to directory")
	clusterFlagSet.StringVar(&cfg.clusterName, "cluster-name", cfg.clusterName, "cluster name")
	kvDefault := ""
	if home := homedir.HomeDir(); home != "" {
		kvDefault = filepath.Join(home, ".holos", "kubeconfig.provisioner")
	}
	kvDefault = getenv("HOLOS_PROVISIONER_KUBECONFIG", kvDefault)
	cfg.kvKubeconfig = kvFlagSet.String("provisioner-kubeconfig", kvDefault, "absolute path to the provisioner cluster kubeconfig file")
	ns := getenv("HOLOS_PROVISIONER_NAMESPACE", DefaultProvisionerNamespace)
	cfg.kvNamespace = kvFlagSet.String("provisioner-namespace", ns, "namespace in the provisioner cluster")
	cfg.txtarIndex = txFlagSet.Int("index", 0, "file number to print if not 0")
	return cfg
}

// Config holds configuration for the whole program, used by main(). The config
// should be initialized early at a well known location in the program lifecycle
// then remain immutable.
type Config struct {
	logConfig            *logger.Config
	writeTo              string
	clusterName          string
	logger               *slog.Logger
	options              *options
	finalized            bool
	writeFlagSet         *flag.FlagSet
	clusterFlagSet       *flag.FlagSet
	kvKubeconfig         *string
	kvNamespace          *string
	kvFlagSet            *flag.FlagSet
	txtarIndex           *int
	txtarFlagSet         *flag.FlagSet
	provisionerClientset kubernetes.Interface
	ServerConfig         *ServerConfig
}

// LogFlagSet returns the logging *flag.FlagSet for use by the command handler.
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

func (c *Config) ServerFlagSet() *flag.FlagSet {
	return c.ServerConfig.FlagSet()
}

// KVFlagSet returns the *flag.FlagSet for kv related commands.
func (c *Config) KVFlagSet() *flag.FlagSet {
	return c.kvFlagSet
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
	return c.logConfig.Vet()
}

// Logger returns a *slog.Logger configured by the user.
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

// Stdin should be used instead of os.Stdin to capture input from tests.
func (c *Config) Stdin() io.Reader {
	return c.options.stdin
}

// Stdout should be used instead of os.Stdout to capture output for tests.
func (c *Config) Stdout() io.Writer {
	return c.options.stdout
}

// Stderr should be used instead of os.Stderr to capture output for tests.
func (c *Config) Stderr() io.Writer {
	return c.options.stderr
}

// WriteTo returns the write to path configured by flags.
func (c *Config) WriteTo() string {
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

// ClusterName returns the cluster name configured by flags.
func (c *Config) ClusterName() string {
	return c.clusterName
}

// KVKubeconfig returns the provisioner cluster kubeconfig path.
func (c *Config) KVKubeconfig() string {
	if c.kvKubeconfig == nil {
		panic("kubeconfig not set")
	}
	return *c.kvKubeconfig
}

// KVNamespace returns the configured namespace to operate against in the provisioner cluster.
func (c *Config) KVNamespace() string {
	if c.kvNamespace == nil {
		return DefaultProvisionerNamespace
	}
	return *c.kvNamespace
}

// TxtarIndex returns the
func (c *Config) TxtarIndex() int {
	if c.txtarIndex == nil {
		return 0
	}
	return *c.txtarIndex
}

// ProvisionerClientset returns a kubernetes client set for the provisioner cluster.
func (c *Config) ProvisionerClientset() (kubernetes.Interface, error) {
	if c.provisionerClientset == nil {
		kcfg, err := clientcmd.BuildConfigFromFlags("", c.KVKubeconfig())
		if err != nil {
			return nil, errors.Wrap(err)
		}
		clientset, err := kubernetes.NewForConfig(kcfg)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		c.provisionerClientset = clientset
	}
	return c.provisionerClientset, nil
}

// getenv is equivalent to os.LookupEnv with a default value.
func getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type ServerConfig struct {
	oidcIssuer     string      // --oidc-issuer
	oidcAudiences  stringSlice // --oidc-audience
	listenAndServe bool        // --no-serve
	listenPort     int         // --listen-port
	metricsPort    int         // --metrics-port
	databaseURI    string
	flagSet        *flag.FlagSet
}

// OIDCIssuer returns the configured oidc issuer url.
func (c *ServerConfig) OIDCIssuer() string {
	return c.oidcIssuer
}

// OIDCAudiences returns the configured allowed id token aud claim values.
func (c *ServerConfig) OIDCAudiences() []string {
	return c.oidcAudiences
}

// DatabaseURI represents the database connection uri.
func (c *ServerConfig) DatabaseURI() string {
	return c.databaseURI
}

// ListenAndServe returns true if the server should listen for and serve requests.
func (c *ServerConfig) ListenAndServe() bool {
	return c.listenAndServe
}

// ListenPort returns the port of the main server.
func (c *ServerConfig) ListenPort() int {
	return c.listenPort
}

// MetricsPort returns the port of the prometheus /metrics scrape endpoint configured by a flag.
func (c *ServerConfig) MetricsPort() int {
	return c.metricsPort
}

func (c *ServerConfig) FlagSet() *flag.FlagSet {
	if c.flagSet != nil {
		return c.flagSet
	}
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.StringVar(&c.oidcIssuer, "oidc-issuer", c.oidcIssuer, "oidc issuer url.")
	f.Var(&c.oidcAudiences, "oidc-audience", "allowed oidc audiences.")
	f.BoolVar(&c.listenAndServe, "serve", true, "listen and serve requests.")
	f.IntVar(&c.listenPort, "listen-port", 3000, "service listen port.")
	f.IntVar(&c.metricsPort, "metrics-port", 9090, "metrics listen port.")
	f.StringVar(&c.databaseURI, "database-url", getenv("DATABASE_URL", c.databaseURI), "database uri (DATABASE_URL)")
	c.flagSet = f
	return f
}

// stringSlice is a comma separated list of string values
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join((*s)[:], ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}
