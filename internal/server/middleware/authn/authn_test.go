package authn_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-jose/go-jose/v3"
	"github.com/holos-run/holos/internal/server/core"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/server/testutils"
	"github.com/stretchr/testify/assert"
)

const issuer = "https://example.com"
const clientID = "holos-cli"

func TestAuthentication(t *testing.T) {
	hf := newHandlerFactory(t)
	t.Run("NewVerifier", func(t *testing.T) {
		// When oidc discovery fails.
		client := &http.Client{Transport: &nullRoundTripper{}}
		log := testutils.TestLogger(t)
		ctx := oidc.ClientContext(logger.NewContext(context.Background(), log), client)
		app := core.AppContext{Context: ctx}.WithLogger(log)
		_, err := authn.NewVerifier(app, issuer)
		assert.Error(t, err)
	})

	t.Run("Handler", func(t *testing.T) {

		t.Run("when_anonymous_deny", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			rw := httptest.NewRecorder()
			hf.NewHandler(t).ServeHTTP(rw, req)
			assert.Equal(t, http.StatusUnauthorized, rw.Result().StatusCode)
		})

		t.Run("when_authorization_no_bearer", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "invalid")
			rw := httptest.NewRecorder()
			hf.NewHandler(t).ServeHTTP(rw, req)
			assert.Equal(t, http.StatusUnauthorized, rw.Result().StatusCode)
		})

		t.Run("when_authorization_bearer_invalid", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer invalid")
			rw := httptest.NewRecorder()
			hf.NewHandler(t).ServeHTTP(rw, req)
			assert.Equal(t, http.StatusUnauthorized, rw.Result().StatusCode)
		})

		t.Run("when_authentic_ok", func(t *testing.T) {
			t.Run("minimal_token", func(t *testing.T) {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+hf.rawIDToken(t))
				rw := httptest.NewRecorder()
				hf.NewHandler(t).ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
			})
		})
	})
}

type signingKey struct {
	keyID string // optional
	priv  interface{}
	pub   interface{}
	alg   jose.SignatureAlgorithm
}

// sign creates a JWS using the private key from the provided payload.
func (s *signingKey) sign(t testing.TB, payload []byte) string {
	privKey := &jose.JSONWebKey{Key: s.priv, Algorithm: string(s.alg), KeyID: s.keyID}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: s.alg, Key: privKey}, nil)
	if err != nil {
		t.Fatal(err)
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}

	data, err := jws.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func (s *signingKey) jwk() jose.JSONWebKey {
	return jose.JSONWebKey{Key: s.pub, Use: "sig", Algorithm: string(s.alg), KeyID: s.keyID}
}

func newRSAKey(t testing.TB) *signingKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1028)
	if err != nil {
		t.Fatal(err)
	}
	return &signingKey{"", privateKey, privateKey.Public(), jose.RS256}
}

type nullRoundTripper struct{}

func (n *nullRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{}
	return response, nil
}

// myHandler extracts the authenticated identity from the context and returns
// values in headers for validation with a response recorder.
func myHandler(w http.ResponseWriter, r *http.Request) {
	id, err := authn.FromContext(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("X-Authn-Issuer", id.Issuer())
	w.Header().Add("X-Authn-Subject", id.Subject())
	w.Header().Add("X-Authn-Name", id.Name())
	w.Header().Add("X-Authn-Email", id.Email())
	if id.Verified() {
		w.Header().Add("X-Authn-Verified", "true")
	} else {
		w.Header().Add("X-Authn-Verified", "false")
	}
}

type handlerFactory struct {
	key      *signingKey
	verifier authn.Verifier
	now      time.Time
}

func (hf *handlerFactory) NewHandler(t testing.TB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.NewContext(r.Context(), testutils.TestLogger(t))
		r = r.WithContext(ctx)
		authn.Handler(hf.verifier, []string{clientID}, http.HandlerFunc(myHandler)).ServeHTTP(w, r)
	})
}

func (hf *handlerFactory) rawIDToken(t testing.TB) string {
	exp := hf.now.Add(time.Hour)
	payload := []byte(fmt.Sprintf(`{
			"iss": "%s",
			"sub": "test_user",
			"aud": "%s",
			"exp": %d
		}`, issuer, clientID, exp.Unix()))
	rawIDToken := hf.key.sign(t, payload)
	return rawIDToken
}

func newHandlerFactory(t testing.TB) *handlerFactory {
	key := newRSAKey(t)
	sks := oidc.StaticKeySet{
		PublicKeys: []crypto.PublicKey{key.jwk().Key},
	}
	now := time.Date(2022, 01, 29, 0, 0, 0, 0, time.UTC)
	config := oidc.Config{
		ClientID: clientID,
		Now:      func() time.Time { return now },
	}
	verifier := oidc.NewVerifier(issuer, &sks, &config)
	return &handlerFactory{key: key, verifier: verifier, now: now}
}
