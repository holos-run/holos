package server_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/holos-run/holos/internal/server"
	"github.com/holos-run/holos/internal/server/app"
	"github.com/holos-run/holos/pkg/version"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	var fastExitError *app.FastExitError

	t.Run("New", func(t *testing.T) {
		cmd := server.New()
		assert.NotNil(t, cmd)
	})

	t.Run("--version", func(t *testing.T) {
		root := newRoot([]string{"--version"})
		err := root.cmd.Execute()
		assert.ErrorAs(t, err, &fastExitError)
		actual := strings.TrimSpace(root.out.String())
		assert.Equal(t, version.Version, actual)
	})

	t.Run("--version-detail", func(t *testing.T) {
		root := newRoot([]string{"--version-detail"})
		err := root.cmd.Execute()
		assert.ErrorAs(t, err, &fastExitError)
		assert.Contains(t, root.out.String(), version.Version)
	})

	// Test as much as possible of the main RunE function.
	t.Run("most_of_run", func(t *testing.T) {
		root := newRoot([]string{"--serve=false"})
		err := root.cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("log_level_invalid", func(t *testing.T) {
		root := newRoot([]string{"--log-level=foo"})
		err := root.cmd.Execute()
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid log level")
		}
	})

	t.Run("log_format_invalid", func(t *testing.T) {
		root := newRoot([]string{"--log-format=foo"})
		err := root.cmd.Execute()
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid log format")
		}
	})

	t.Run("issuer_invalid", func(t *testing.T) {
		expected := "oidc issuer must start with https://"
		root := newRoot([]string{"--oidc-issuer=foo"})
		err := root.cmd.Execute()
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), expected)
		}
	})

	t.Run("log_levels", func(t *testing.T) {
		var validLogLevels = []string{"debug", "info", "warn", "error"}
		for _, level := range validLogLevels {
			t.Run(level, func(t *testing.T) {
				root := newRoot([]string{"--serve=false", "--log-level", level})
				err := root.cmd.Execute()
				assert.NoError(t, err)
			})
		}
	})

	t.Run("log_formats", func(t *testing.T) {
		var validLogFormats = []string{"text", "json"}
		for _, format := range validLogFormats {
			t.Run(format, func(t *testing.T) {
				root := newRoot([]string{"--serve=false", "--log-format", format})
				err := root.cmd.Execute()
				assert.NoError(t, err)
			})
		}
	})
}

type oidcMocker struct{}

func (n *oidcMocker) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     map[string][]string{},
	}
	switch path := req.URL.Path; path {
	case "/.well-known/openid-configuration":
		s := `{
			"issuer": "https://example.com",
			"authorization_endpoint": "https://example.com/auth",
			"token_endpoint": "https://example.com/token",
			"jwks_uri": "https://example.com/keys",
			"userinfo_endpoint": "https://example.com/userinfo"
        }`
		response.Header.Set("Content-Type", "application/json")
		response.Body = io.NopCloser(strings.NewReader(s))
	default:
		response.StatusCode = http.StatusNotFound
	}

	return response, nil
}

// mockIdentityProvider mocks oidc discovery responses from an issuer.
func mockIdentityProvider(ctx context.Context) context.Context {
	client := &http.Client{Transport: &oidcMocker{}}
	return oidc.ClientContext(ctx, client)
}

// root is a test harness for the root command and output.
type root struct {
	cmd *cobra.Command
	err *bytes.Buffer
	out *bytes.Buffer
}

func newRoot(args []string) root {
	rc := server.New(app.WithIssuer("https://example.com"))
	r := root{
		cmd: rc,
		err: new(bytes.Buffer),
		out: new(bytes.Buffer),
	}
	r.cmd.SetArgs(args)
	r.cmd.SetErr(r.err)
	r.cmd.SetOut(r.out)
	// mock http client for the coreos/oidc module
	ctx := mockIdentityProvider(context.Background())
	r.cmd.SetContext(ctx)
	return r
}
