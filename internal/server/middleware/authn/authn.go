// Package authn provides the middleware handler responsible for authenticating
// requests and adding the Identity to the request context. @todo rename this
// package to authn (authentication) to distinguish from authz (authorization)
package authn

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
)

const Header = "x-oidc-id-token"

// Verifier is the interface that wraps the basic Verify method to verify an
// oidc id token is authentic. Intended for use in request handlers.
type Verifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

// Identity is the interface that defines an authenticated subject (principal,
// person or service) in the system. The methods correspond to oidc claims for
// the cli api client using scopes of, "email profile groups offline_access"
//
// The primary use case is Dex connected to Google using the Google connector
// with a groups reader service account to fetch group membership.
//
// Behavior to keep in mind with Dex v2.37.0 and the `google` connector:
//
// 1. There is only one refresh token stored for each user/client pair. 2. Dex
// does not return the `name` claim in the id token returned from exchanging a
// refresh token. Google specifies they may omit the name claim. oauth spec says
// providers may omit the name in refresh responses.
type Identity interface {
	// Issuer is the oidc issuer url.
	Issuer() string
	// Subject is the unique id of the user within the context of the issuer.
	Subject() string
	// Email address of the user.
	Email() string
	// Verified is true if the email address has been verified by the identity provider.
	Verified() bool
	// Name is usually set on the initial id token, often omitted by google in refreshed id tokens.
	Name() string
}

// key is an unexported type for keys defined in this package to prevent
// collisions with keys defined in other packages.
type key int

// https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/context/context.go;l=140-158
// userKey is the key for Identity providers in Contexts.  It is unexported,
// clients use NewContext and FromContext instead of this key directly.
var userKey key

// NewContext returns a new Context that carries value u. Use FromContext
// to retrieve the value.
func NewContext(ctx context.Context, u Identity) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the value previously stored in ctx by NewContext or nil.
func FromContext(ctx context.Context) (Identity, error) {
	if user, ok := ctx.Value(userKey).(Identity); ok {
		return user, nil
	}
	return nil, errors.New("no user in request context")
}

// NewVerifier returns an *oidc.IDTokenVerifier that implements Verifier from an
// oidc.Provider for issuer which performs jwks .well-known discovery.
func NewVerifier(ctx context.Context, log *slog.Logger, issuer string) (*oidc.IDTokenVerifier, error) {
	var err error
	var oidcProvider *oidc.Provider
	for i := 1; i < 30; i++ {
		oidcProvider, err = oidc.NewProvider(ctx, issuer)
		if err != nil {
			if strings.Contains(err.Error(), "connect: connection refused") {
				log.DebugContext(ctx, "could not get oidc provider, the service mesh sidecar or network may not be ready, retrying", "err", err, "try", i, "max", 30)
				time.Sleep(1 * time.Second)
			} else {
				log.DebugContext(ctx, "could not get oidc provider", "err", err)
				break
			}
		} else {
			if i > 1 {
				log.DebugContext(ctx, "ok: got oidc provider", "try", i)
			}
			break
		}
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}
	// We allow tokens from multiple client ids (web, cli), they are checked in the handler.
	return oidcProvider.Verifier(&oidc.Config{SkipClientIDCheck: true}), nil
}

type claims struct {
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Name     string `json:"name"`
}

type user struct {
	claims claims
}

func (u user) Issuer() string {
	return u.claims.Issuer
}

func (u user) Subject() string {
	return u.claims.Subject
}

func (u user) Name() string {
	return u.claims.Name
}

func (u user) Email() string {
	return u.claims.Email
}

func (u user) Verified() bool {
	return u.claims.Verified
}

// Handler returns a handler that verifies the request is authentic and adds a
// Identity to the request context.
func Handler(v Verifier, allowedAudiences []string, header string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rawIDToken string
		start := time.Now()
		// Acquire the logger
		log := logger.FromContext(r.Context()).With("handler", "auth")

		// Check the X-Auth-Request-Access-Token header set by Istio ExternalAuthorization
		if rawIDToken == "" {
			rawIDToken = r.Header.Get(header)
		}

		// Validate the authorization bearer token
		if rawIDToken == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				reason := "request missing authorization header"
				log.Debug("could not authenticate", "reason", reason)
				http.Error(w, reason, http.StatusUnauthorized)
				return
			}
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) != 2 {
				reason := "request authorization header is not a bearer token"
				log.Debug("could not authenticate", "reason", reason)
				http.Error(w, reason, http.StatusUnauthorized)
				return
			}
			rawIDToken = splitToken[1]
		}

		idToken, err := v.Verify(r.Context(), rawIDToken)
		if err != nil {
			log.Error("invalid authorization bearer id token", "err", err)
			http.Error(w, "invalid authorization bearer id token", http.StatusUnauthorized)
			return
		}

		// Check audiences
		var audOK bool
		for _, expectedAud := range allowedAudiences {
			for _, haveAud := range idToken.Audience {
				if haveAud == expectedAud {
					audOK = true
					break
				}
			}
			if audOK {
				break
			}
		}

		if !audOK {
			log.Error("audience not allowed", "expected", allowedAudiences, "got", idToken.Audience)
			http.Error(w, "audience not allowed", http.StatusUnauthorized)
			return
		}

		// ID Token is valid, extract the claims into the user struct
		u := user{}
		if err = idToken.Claims(&u.claims); err != nil {
			log.Error("could not extract claims from id token", "err", err)
			http.Error(w, "could not extract claims from id token", http.StatusInternalServerError)
			return
		}

		// Add the user to the context and update the logger in the context
		ctx := r.Context()
		// Log only the subject to protect pii (email, name)
		userLogger := logger.FromContext(ctx).With("sub", u.Subject())
		userCtx := NewContext(logger.NewContext(ctx, userLogger), u)

		next.ServeHTTP(w, r.WithContext(userCtx))
		userLogger.DebugContext(ctx, "response", "duration", time.Since(start))
	})
}
