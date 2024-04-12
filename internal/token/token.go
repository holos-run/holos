// Package token obtains, caches, and provides an ID token to authenticate to the holos api server.
package token

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/holos-run/holos/pkg/errors"
	"github.com/int128/kubelogin/pkg/infrastructure/browser"
	"github.com/int128/kubelogin/pkg/infrastructure/clock"
	"github.com/int128/kubelogin/pkg/infrastructure/logger"
	"github.com/int128/kubelogin/pkg/oidc"
	"github.com/int128/kubelogin/pkg/oidc/client"
	"github.com/int128/kubelogin/pkg/tlsclientconfig"
	"github.com/int128/kubelogin/pkg/tlsclientconfig/loader"
	"github.com/int128/kubelogin/pkg/tokencache"
	"github.com/int128/kubelogin/pkg/tokencache/repository"
	"github.com/int128/kubelogin/pkg/usecases/authentication"
	"github.com/int128/kubelogin/pkg/usecases/authentication/authcode"
	"github.com/int128/kubelogin/pkg/usecases/authentication/devicecode"
	"github.com/int128/kubelogin/pkg/usecases/authentication/ropc"
	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

var CacheDir = expandHomedir(filepath.Join("~", ".holos", "cache"))

// Token represents an authorization bearer token. Token is useful as an output
// dto of the Tokener service use case.
type Token struct {
	// Bearer is the oidc token for the authorization: bearer header
	Bearer string
	// Expiry is the expiration time of the id token
	Expiry time.Time
	// Pretty is the JSON encoding of the token claims
	Pretty string
	// claims represent decoded claims
	claims *Claims
}

func (t Token) Claims() *Claims {
	if t.claims == nil {
		json.Unmarshal([]byte(t.Pretty), &t.claims)
	}
	return t.claims
}

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// NewConfig returns a Config with default values.
func NewConfig() Config {
	return Config{
		Issuer:      "https://login.ois.run",
		ClientID:    "262479925313799528@holos_platform",
		Scopes:      []string{"openid", "email", "profile", "groups", "offline_access"},
		ExtraScopes: []string{"urn:zitadel:iam:org:domain:primary:openinfrastructure.co"},
	}
}

type Config struct {
	Issuer       string
	ClientID     string
	Scopes       stringSlice
	ExtraScopes  stringSlice
	ForceRefresh bool
	flagSet      *flag.FlagSet
}

func (c *Config) FlagSet() *flag.FlagSet {
	if c.flagSet != nil {
		return c.flagSet
	}
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&c.Issuer, "oidc-issuer", c.Issuer, "oidc token issuer url.")
	flags.StringVar(&c.ClientID, "oidc-client-id", c.ClientID, "oidc client id.")
	flags.Var(&c.Scopes, "oidc-scopes", "required oidc scopes")
	flags.Var(&c.ExtraScopes, "oidc-extra-scopes", "optional oidc scopes")
	flags.BoolVar(&c.ForceRefresh, "oidc-force-refresh", c.ForceRefresh, "force refresh")
	c.flagSet = flags
	return flags
}

// Get returns an oidc token for use as an authorization bearer http header.
func Get(ctx context.Context, log *slog.Logger, cfg Config) (*Token, error) {
	var scopes []string
	scopes = append(scopes, cfg.Scopes...)
	scopes = append(scopes, cfg.ExtraScopes...)
	provider := oidc.Provider{
		IssuerURL:   cfg.Issuer,
		ClientID:    cfg.ClientID,
		UsePKCE:     true,
		ExtraScopes: scopes,
	}

	authenticationOptions := authenticationOptions{
		GrantType:                   "auto",
		ListenAddress:               defaultListenAddress,
		AuthenticationTimeoutSec:    180,
		RedirectURLHostname:         "localhost",
		RedirectURLAuthCodeKeyboard: oobRedirectURI,
	}

	grantOptionSet, err := authenticationOptions.grantOptionSet()
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not login: %w", err))
	}

	tlsClientConfig := tlsclientconfig.Config{}

	tokenCacheKey := tokencache.Key{
		IssuerURL:      provider.IssuerURL,
		ClientID:       provider.ClientID,
		ClientSecret:   provider.ClientSecret,
		ExtraScopes:    provider.ExtraScopes,
		CACertFilename: strings.Join(tlsClientConfig.CACertFilename, ","),
		CACertData:     strings.Join(tlsClientConfig.CACertData, ","),
		SkipTLSVerify:  tlsClientConfig.SkipTLSVerify,
	}

	if grantOptionSet.ROPCOption != nil {
		tokenCacheKey.Username = grantOptionSet.ROPCOption.Username
	}

	tokenCacheRepository := &repository.Repository{}

	cachedTokenSet, err := tokenCacheRepository.FindByKey(CacheDir, tokenCacheKey)
	if err != nil {
		slog.Debug("could not find a token cache (continuing)", "err", err, "handled", true)
	}

	// Construct input for the Authentication service use case
	authenticationInput := authentication.Input{
		Provider:        provider,
		GrantOptionSet:  grantOptionSet,
		CachedTokenSet:  cachedTokenSet,
		TLSClientConfig: tlsClientConfig,
		ForceRefresh:    cfg.ForceRefresh,
	}

	var slogger logger.Interface = &holosLogger{log: log}

	clock := &clock.Real{}

	auth := &authentication.Authentication{
		ClientFactory: &client.Factory{
			Loader: loader.Loader{},
			Clock:  clock,
			Logger: slogger,
		},
		Logger: slogger,
		Clock:  clock,
		AuthCodeBrowser: &authcode.Browser{
			Browser: &browser.Browser{},
			Logger:  slogger,
		},
	}

	authenticationOutput, err := auth.Do(ctx, authenticationInput)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	idTokenClaims, err := authenticationOutput.TokenSet.DecodeWithoutVerify()
	if err != nil {
		slog.Debug("could not get token claims", "err", err, "handled", false)
		return nil, fmt.Errorf("could not get token claims: %w", err)
	}

	if authenticationOutput.AlreadyHasValidIDToken {
		slog.Debug("existing token valid", "refreshed", 0, "exp", idTokenClaims.Expiry)
	} else {
		slog.Debug("new token valid", "refreshed", 1, "exp", idTokenClaims.Expiry)
		if err := tokenCacheRepository.Save(CacheDir, tokenCacheKey, authenticationOutput.TokenSet); err != nil {
			slog.Debug("could not save token cache", "err", err, "handled", 0)
			return nil, fmt.Errorf("could not save token cache: %w", err)
		}
	}

	token := &Token{
		Bearer: authenticationOutput.TokenSet.IDToken,
		Expiry: idTokenClaims.Expiry,
		Pretty: idTokenClaims.Pretty,
	}
	return token, nil
}

var defaultListenAddress = []string{"127.0.0.1:8000", "127.0.0.1:18000"}
var allGrantType = strings.Join([]string{
	"auto",
	"authcode",
	"authcode-keyboard",
	"password",
	"device-code",
}, "|")

const oobRedirectURI = "urn:ietf:wg:oauth:2.0:oob"

func expandHomedir(s string) string {
	if !strings.HasPrefix(s, "~") {
		return s
	}
	return filepath.Join(homedir.HomeDir(), strings.TrimPrefix(s, "~"))
}

type authenticationOptions struct {
	GrantType                   string
	ListenAddress               []string
	AuthenticationTimeoutSec    int
	SkipOpenBrowser             bool
	BrowserCommand              string
	LocalServerCertFile         string
	LocalServerKeyFile          string
	OpenURLAfterAuthentication  string
	RedirectURLHostname         string
	RedirectURLAuthCodeKeyboard string
	AuthRequestExtraParams      map[string]string
	Username                    string
	Password                    string
}

func (o *authenticationOptions) grantOptionSet() (s authentication.GrantOptionSet, err error) {
	switch {
	case o.GrantType == "authcode" || (o.GrantType == "auto" && o.Username == ""):
		s.AuthCodeBrowserOption = &authcode.BrowserOption{
			BindAddress:                o.ListenAddress,
			SkipOpenBrowser:            o.SkipOpenBrowser,
			BrowserCommand:             o.BrowserCommand,
			AuthenticationTimeout:      time.Duration(o.AuthenticationTimeoutSec) * time.Second,
			LocalServerCertFile:        o.LocalServerCertFile,
			LocalServerKeyFile:         o.LocalServerKeyFile,
			OpenURLAfterAuthentication: o.OpenURLAfterAuthentication,
			RedirectURLHostname:        o.RedirectURLHostname,
			AuthRequestExtraParams:     o.AuthRequestExtraParams,
		}
	case o.GrantType == "authcode-keyboard":
		s.AuthCodeKeyboardOption = &authcode.KeyboardOption{
			AuthRequestExtraParams: o.AuthRequestExtraParams,
			RedirectURL:            o.RedirectURLAuthCodeKeyboard,
		}
	case o.GrantType == "password" || (o.GrantType == "auto" && o.Username != ""):
		s.ROPCOption = &ropc.Option{
			Username: o.Username,
			Password: o.Password,
		}
	case o.GrantType == "device-code":
		s.DeviceCodeOption = &devicecode.Option{
			SkipOpenBrowser: o.SkipOpenBrowser,
			BrowserCommand:  o.BrowserCommand,
		}
	default:
		err = fmt.Errorf("grant-type must be one of (%s)", allGrantType)
	}
	return
}

// holosLogger implements the int128/kubelogin logger.Interface
type holosLogger struct {
	log *slog.Logger
}

func (*holosLogger) AddFlags(f *pflag.FlagSet) {}
func (l *holosLogger) Printf(format string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(format, args...))
}
func (l *holosLogger) Infof(format string, args ...interface{}) {
	l.Printf(format, args...)
}
func (l *holosLogger) V(level int) logger.Verbose {
	return l
}
func (*holosLogger) IsEnabled(level int) bool {
	return true
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
