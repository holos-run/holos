package server

import (
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/holos-run/holos/internal/server/app"
	"github.com/holos-run/holos/internal/server/ent"
	"github.com/holos-run/holos/internal/server/frontend"
	"github.com/holos-run/holos/internal/server/handler"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/server/service/gen/holos/v1alpha1/holosconnect"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// @title Holos Server
// @version 0.1
// @description Holos is a holistically integrated software development platform.

// @contact.name Open Infrastructure Services LLC
// @contact.url https://openinfrastructure.co

// @license.name TBD
// @license.url https://openinfrastructure.co

// @host localhost:8443
// @BasePath /
// @schemes https

var (
	healthy int32
	ready   int32
)

type Config struct {
	HttpServerTimeout     time.Duration
	ServerShutdownTimeout time.Duration
	Host                  string
	Port                  string
	Unhealthy             bool
	Unready               bool
	MetricsPort           int
	OIDCIssuer            string
	OIDCAudiences         []string
}

type Server struct {
	app           app.App
	mux           *http.ServeMux
	handler       http.Handler
	config        *Config
	db            *ent.Client
	authenticator authn.Verifier
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

type Middleware func(next http.Handler) http.Handler

func NewServer(app app.App, config *Config, db *ent.Client, verifier authn.Verifier) (*Server, error) {
	mux := http.NewServeMux()
	srv := &Server{
		app:           app,
		mux:           mux,
		handler:       h2c.NewHandler(mux, &http2.Server{}),
		config:        config,
		db:            db,
		authenticator: verifier,
	}

	srv.registerHandlers()
	if err := srv.registerConnectRpc(); err != nil {
		return srv, wrapper.Wrap(err)
	}

	return srv, nil
}

// middlewares wraps the handler with our standard middleware chain.
func (s *Server) middlewares(handler http.Handler) http.Handler {
	return logger.LoggingMiddleware(s.app.Logger)(s.uiMiddleware(handler))
}

func (s *Server) notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := s.app.Logger
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404 page not found\n")); err != nil {
			log.ErrorContext(ctx, "could not write", "err", err)
		}
	})
}

func (s *Server) registerHandlers() {
	// Prometheus metrics
	s.mux.Handle("/metrics", promhttp.Handler())
	// Main entrypoint for the frontend interface
	fsHandler := http.FileServer(http.FS(frontend.Root()))
	s.mux.Handle(frontend.Path, s.middlewares(frontend.SPAFileServer(fsHandler)))
	// Handle 404 errors server side for static paths like assets and logos.
	staticHandler := s.middlewares(fsHandler)
	s.mux.Handle(frontend.Path+"assets/", staticHandler)
	s.mux.Handle(frontend.Path+"logos/", staticHandler)

	// Redirect GET on / to the frontend
	s.mux.Handle("/", s.middlewares(s.notFoundHandler()))
}

func (s *Server) registerConnectRpc() error {
	// Validator for all rpc messages
	validator, err := validate.NewInterceptor()
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not initialize proto validation interceptor: %w", err))
	}

	h := handler.NewHolosHandler(s.app, s.db)
	holosPath, holosHandler := holosconnect.NewHolosServiceHandler(h, connect.WithInterceptors(validator))
	authenticatingHandler := authn.Handler(s.authenticator, s.config.OIDCAudiences, holosHandler)
	s.mux.Handle(holosPath, s.middlewares(authenticatingHandler))

	return nil
}

// uiMiddleware redirects GET / to frontend.Path to load the web ui
func (s *Server) uiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			http.Redirect(w, r, frontend.Path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) ListenAndServe() (*http.Server, *int32, *int32) {
	go s.startMetricsServer()

	// create the http server
	srv := s.startServer()

	// signal Kubernetes the server is ready to receive traffic
	if !s.config.Unhealthy {
		atomic.StoreInt32(&healthy, 1)
	}
	if !s.config.Unready {
		atomic.StoreInt32(&ready, 1)
	}

	return srv, &healthy, &ready
}

func (s *Server) startServer() *http.Server {
	// determine if the port is specified
	if s.config.Port == "0" {
		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:         s.config.Host + ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	httpLog := s.app.Logger.With("addr", srv.Addr, "server", "http")

	// start the server in the background
	go func() {
		httpLog.InfoContext(s.app.Context, "listening for http requests")
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			httpLog.DebugContext(s.app.Context, "server closed", "lifecycle", "end")
		} else {
			httpLog.ErrorContext(s.app.Context, "could not listen for http requests", "err", err, "lifecycle", "end", "exit", 2)
			os.Exit(2)
		}
	}()

	// return the server and routine
	return srv
}

func (s *Server) startMetricsServer() {
	ctx, log := s.app.ContextLogger()
	if s.config.MetricsPort < 1 {
		log.WarnContext(ctx, "metrics disabled", "suggestion", "enable with flag --metrics-port=9090")
	}
	mux := http.DefaultServeMux
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.ErrorContext(ctx, "could not write", "err", err)
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.config.MetricsPort),
		Handler: mux,
	}

	srvLog := log.With("addr", srv.Addr, "server", "metrics")
	srvLog.InfoContext(ctx, "listening for prom requests")
	// Blocks the go routine, always returns a non-nil error
	if err := srv.ListenAndServe(); err != nil {
		srvLog.ErrorContext(ctx, "could not listen for prom requests", "err", err)
	}
}

type ArrayResponse []string
type MapResponse map[string]string
