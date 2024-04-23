package server

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/frontend"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/handler"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/service/gen/holos/v1alpha1/holosconnect"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	healthy int32
	ready   int32
)

type Server struct {
	cfg           *holos.Config
	db            *ent.Client
	mux           *http.ServeMux
	handler       http.Handler
	authenticator authn.Verifier
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

type Middleware func(next http.Handler) http.Handler

func NewServer(cfg *holos.Config, db *ent.Client, verifier authn.Verifier) (*Server, error) {
	mux := http.NewServeMux()
	srv := &Server{
		cfg:           cfg,
		db:            db,
		mux:           mux,
		handler:       h2c.NewHandler(mux, &http2.Server{}),
		authenticator: verifier,
	}

	srv.registerHandlers()
	if err := srv.registerConnectRpc(); err != nil {
		return srv, errors.Wrap(err)
	}

	return srv, nil
}

// middlewares wraps the handler with our standard middleware chain.
func (s *Server) middlewares(handler http.Handler) http.Handler {
	return logger.LoggingMiddleware(s.cfg.Logger())(s.uiMiddleware(handler))
}

func (s *Server) notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := s.cfg.Logger()
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404 page not found\n")); err != nil {
			log.ErrorContext(ctx, "could not write", "err", err)
		}
	})
}

func (s *Server) registerHandlers() {
	// Prometheus metrics
	s.mux.Handle("/metrics", promhttp.Handler())
	// Main entrypoint for the frontend interface.
	spaFileServer := frontend.NewSPAFileServer(s.cfg.ServerConfig.OIDCIssuer())
	fsHandler := http.FileServer(http.FS(frontend.Root()))
	s.mux.Handle(frontend.Path, s.middlewares(spaFileServer(fsHandler)))
	// Handle 404 errors server side for static paths like assets and logos.
	staticHandler := s.middlewares(fsHandler)
	s.mux.Handle(frontend.Path+"assets/", staticHandler)
	s.mux.Handle(frontend.Path+"logos/", staticHandler)

	// Redirect GET on / to the frontend
	s.mux.Handle("/", s.middlewares(s.notFoundHandler()))
}

// handle wraps handler with holos authentication and registers the handler with the server mux.
func (s *Server) handle(pattern string, handler http.Handler) {
	authenticatingHandler := authn.Handler(
		s.authenticator,
		s.cfg.ServerConfig.OIDCAudiences(),
		s.cfg.ServerConfig.AuthHeader(),
		handler,
	)
	s.mux.Handle(pattern, authenticatingHandler)
}

func (s *Server) registerConnectRpc() error {
	// Validator for all rpc messages
	validator, err := validate.NewInterceptor()
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not initialize proto validation interceptor: %w", err))
	}

	opts := connect.WithInterceptors(validator)

	s.handle(holosconnect.NewUserServiceHandler(handler.NewUserHandler(s.db), opts))
	s.handle(holosconnect.NewOrganizationServiceHandler(handler.NewOrganizationHandler(s.db), opts))

	reflector := grpcreflect.NewStaticReflector(
		holosconnect.UserServiceName,
		holosconnect.OrganizationServiceName,
	)

	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

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
	atomic.StoreInt32(&healthy, 1)
	atomic.StoreInt32(&ready, 1)

	return srv, &healthy, &ready
}

func (s *Server) startServer() *http.Server {
	// determine if the port is specified
	if s.cfg.ServerConfig.ListenPort() == 0 {
		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.ServerConfig.ListenPort()),
		Handler: s.handler,
		// WriteTimeout: s.cfg.ServerConfig.HttpServerTimeout,
		// ReadTimeout:  s.cfg.ServerConfig.HttpServerTimeout,
		// IdleTimeout:  2 * s.cfg.ServerConfig.HttpServerTimeout,
	}

	httpLog := s.cfg.Logger().With("addr", srv.Addr, "server", "http")

	// start the server in the background
	go func() {
		httpLog.Info("listening for http requests")
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			httpLog.Debug("server closed", "lifecycle", "end")
		} else {
			httpLog.Error("could not listen for http requests", "err", err, "lifecycle", "end", "exit", 2)
			os.Exit(2)
		}
	}()

	// return the server and routine
	return srv
}

func (s *Server) startMetricsServer() {
	log := s.cfg.Logger()
	if s.cfg.ServerConfig.MetricsPort() < 1 {
		log.Warn("metrics disabled", "suggestion", "enable with flag --metrics-port=9090")
	}
	mux := http.DefaultServeMux
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Error("could not write", "err", err)
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.ServerConfig.MetricsPort()),
		Handler: mux,
	}

	srvLog := log.With("addr", srv.Addr, "server", "metrics")
	srvLog.Info("listening for prom requests")
	// Blocks the go routine, always returns a non-nil error
	if err := srv.ListenAndServe(); err != nil {
		srvLog.Error("could not listen for prom requests", "err", err)
	}
}

type ArrayResponse []string
type MapResponse map[string]string
