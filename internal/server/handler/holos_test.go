package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/google/uuid"
	"github.com/holos-run/holos/internal/server/db"
	"github.com/holos-run/holos/internal/server/ent"
	"github.com/holos-run/holos/internal/server/handler"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/testutils"
	"github.com/holos-run/holos/pkg/holos"
	holosSvc "github.com/holos-run/holos/service/gen/holos/v1alpha1"
	"github.com/holos-run/holos/service/gen/holos/v1alpha1/holosconnect"
	"github.com/stretchr/testify/assert"
)

const (
	authHeader = "Authorized-User"
)

func TestHolosService(t *testing.T) {
	t.Parallel()
	testUserCreateUpdate := func(t *testing.T, client holosconnect.HolosServiceClient) {
		t.Run("CreateUser", func(t *testing.T) {
			t.Run("registration", func(t *testing.T) {
				// The id token claims should be sufficient for user registration. Authenticated
				// clients need only pass in the zero value to register themselves.
				request, authnIdentity := newAuthenticRequest(&holosSvc.RegisterUserRequest{})
				response, err := client.RegisterUser(testutils.LogCtx(t), request)
				if err != nil {
					t.Errorf("unexpected: %v", err)
				} else {
					actual := response.Msg.GetUser()
					assert.Equal(t, authnIdentity.Name(), actual.GetName(), "name does not match auth claims")
					assert.Equal(t, authnIdentity.Email(), actual.GetEmail(), "email does not match auth claims")
				}
			})
			t.Run("name", func(t *testing.T) {
				// An authenticated client can set their name.
				expected := "Bob"
				req, authnIdentity := newAuthenticRequest(&holosSvc.RegisterUserRequest{
					Name: &expected,
				})
				resp, err := client.RegisterUser(testutils.LogCtx(t), req)
				if err != nil {
					t.Errorf("unexpected: %v", err)
				} else {
					actual := resp.Msg.GetUser()
					assert.Equal(t, expected, actual.GetName(), "name returned does not match request")
					assert.Equal(t, authnIdentity.Email(), actual.GetEmail(), "email does not match auth claims")
				}
			})
			t.Run("name_length_limit", func(t *testing.T) {
				// Name cannot be longer than 100 characters.
				name := strings.Repeat("X", 101)
				req, _ := newAuthenticRequest(&holosSvc.RegisterUserRequest{
					Name: &name,
				})
				_, err := client.RegisterUser(testutils.LogCtx(t), req)
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), "name: value length must be at most 100 characters")
				}
			})
		})
	}

	testMatrix := func(t *testing.T, server *httptest.Server) {
		run := func(t *testing.T, opts ...connect.ClientOption) {
			t.Helper()
			client := holosconnect.NewHolosServiceClient(server.Client(), server.URL, opts...)
			testUserCreateUpdate(t, client)
		}
		t.Run("connect", func(t *testing.T) {
			t.Run("proto", func(t *testing.T) {
				run(t)
			})
			t.Run("proto_gzip", func(t *testing.T) {
				run(t, connect.WithSendGzip())
			})
			t.Run("json_gzip", func(t *testing.T) {
				run(
					t,
					connect.WithProtoJSON(),
					connect.WithSendGzip(),
				)
			})
		})
		t.Run("grpc", func(t *testing.T) {
			t.Run("proto", func(t *testing.T) {
				run(t, connect.WithGRPC())
			})
			t.Run("proto_gzip", func(t *testing.T) {
				run(t, connect.WithGRPC(), connect.WithSendGzip())
			})
			t.Run("json_gzip", func(t *testing.T) {
				run(
					t,
					connect.WithGRPC(),
					connect.WithProtoJSON(),
					connect.WithSendGzip(),
				)
			})
		})
		t.Run("grpcweb", func(t *testing.T) {
			t.Run("proto", func(t *testing.T) {
				run(t, connect.WithGRPCWeb())
			})
			t.Run("proto_gzip", func(t *testing.T) {
				run(t, connect.WithGRPCWeb(), connect.WithSendGzip())
			})
			t.Run("json_gzip", func(t *testing.T) {
				run(
					t,
					connect.WithGRPCWeb(),
					connect.WithProtoJSON(),
					connect.WithSendGzip(),
				)
			})
		})
	}

	requestMessageValidator, err := validate.NewInterceptor()
	if err != nil {
		panic(err)
	}
	routePath, holosHandler := holosconnect.NewHolosServiceHandler(
		// Inject the logger here.
		handler.NewHolosHandler(newDatabaseClient(t)),
		connect.WithInterceptors(requestMessageValidator),
	)

	// Add a fake authenticated user to the request context.
	mux := http.NewServeMux()
	mux.Handle(routePath, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		holosHandler.ServeHTTP(resp, withAuthContext(t, req))
	}))

	// Run the test matrix suites
	t.Run("http1", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(mux)
		t.Cleanup(server.Close)
		testMatrix(t, server)
	})
	t.Run("http2", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewUnstartedServer(mux)
		server.EnableHTTP2 = true
		server.StartTLS()
		t.Cleanup(server.Close)
		testMatrix(t, server)
	})
}

type claims struct {
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Name     string `json:"name"`
}

// user mocks the authn.Identity stored in the request context by the authn middleware.
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

// newDatabaseClient returns a new database client for testing.
func newDatabaseClient(t *testing.T) *ent.Client {
	cfg := holos.New(holos.Logger(testutils.TestLogger(t)))
	// Connect to the database
	var dbf db.ClientFactory = db.NewMemoryClientFactory(cfg)
	conn, err := dbf.New()
	dbClient := conn.Client
	if err != nil {
		panic(err)
	}
	// Automatic migration
	if err = dbClient.Schema.Create(testutils.LogCtx(t)); err != nil {
		panic(err)
	}
	return dbClient
}

func newUser() user {
	sub := uuid.New().String()
	usr := user{
		claims: claims{
			Issuer:   "https://example.com",
			Subject:  sub,
			Name:     "Alice",
			Email:    fmt.Sprintf("alice-%s@example.com", sub),
			Verified: true,
		},
	}
	return usr
}

// withAuthHeader adds fake authn user info to the request headers to mock
// authorization bearer tokens.
func withAuthHeader[T any](request *connect.Request[T]) user {
	usr := newUser()
	userJsonBytes, err := json.Marshal(usr.claims)
	if err != nil {
		panic(err)
	}
	request.Header().Set(authHeader, string(userJsonBytes))
	return usr
}

// newAuthenticRequest returns a request with valid test authorization headers
// for user.
func newAuthenticRequest[T any](message *T) (*connect.Request[T], user) {
	request := connect.NewRequest(message)
	userIdentity := withAuthHeader(request)
	return request, userIdentity
}

// withAuthContext copies fake authn user info from request headers to the
// request context using the same authn.Identity the production authentication
// decorator uses.
func withAuthContext(t *testing.T, request *http.Request) *http.Request {
	var usr user
	err := json.Unmarshal([]byte(request.Header.Get(authHeader)), &usr.claims)
	if assert.NoError(t, err) {
		userCtx := authn.NewContext(request.Context(), usr)
		return request.WithContext(userCtx)
	}
	// panic maybe?
	return request
}
