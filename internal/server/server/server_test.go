package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/frontend"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/db"
	"github.com/holos-run/holos/internal/server/server"
	"github.com/holos-run/holos/internal/server/testutils"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	cfg := holos.New(holos.Logger(testutils.TestLogger(t)))
	client := newClient(cfg)
	srv, err := server.NewServer(cfg, client, &fakeVerifier{})
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	t.Run("when_browser_reloads_spa_managed_route", func(t *testing.T) {
		t.Run("returns_app_index_for_ui_profile", func(t *testing.T) {
			rr := httptest.NewRecorder()
			srv.Mux().ServeHTTP(rr, newReq(frontend.Path+"profile"))
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode, "should be 200 OK")
			contentType := rr.Result().Header.Get("Content-Type")
			assert.Equal(t, "text/html; charset=utf-8", contentType, "should be text/html")
		})
		t.Run("returns_404_for_ui_assets_does_not_exist", func(t *testing.T) {
			rr := httptest.NewRecorder()
			srv.Mux().ServeHTTP(rr, newReq("/ui/assets/does-not-exist"))
			assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode, "should exclude static paths")
		})
		t.Run("returns_404_for_ui_logos_does_not_exist", func(t *testing.T) {
			rr := httptest.NewRecorder()
			srv.Mux().ServeHTTP(rr, newReq("/ui/logos/does-not-exist"))
			assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode, "should exclude static paths")
		})
	})
}

// fakeVerifier implements the Verifier interface and returns fake IDTokens for
// testing.
type fakeVerifier struct{}

func (p fakeVerifier) Verify(context.Context, string) (*oidc.IDToken, error) {
	now := time.Now()
	idToken := &oidc.IDToken{
		Issuer:          "https://example.com",
		Audience:        []string{"example"},
		Subject:         "18be380c-48a0-48ad-8c9f-f873733f24be",
		Expiry:          now.Add(time.Duration(-30) * time.Minute),
		IssuedAt:        now.Add(time.Duration(30) * time.Minute),
		Nonce:           "",
		AccessTokenHash: "",
	}
	return idToken, nil
}
func newClient(cfg *holos.Config) *ent.Client {
	// Connect to the database
	var dbf db.ClientFactory = db.NewMemoryClientFactory(cfg)
	conn, err := dbf.New()
	dbClient := conn.Client
	if err != nil {
		panic(err)
	}
	// Automatic migration
	if err = dbClient.Schema.Create(context.Background()); err != nil {
		panic(err)
	}
	return dbClient
}

func newReq(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return req
}
