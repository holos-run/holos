package server

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/lmittmann/tint"
	"github.com/sethvargo/go-retry"
	"github.com/spf13/cobra"

	"github.com/holos-run/holos/internal/server/app"
	"github.com/holos-run/holos/internal/server/db"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/server/server"
	"github.com/holos-run/holos/internal/server/signals"
	"github.com/holos-run/holos/pkg/version"
	"github.com/holos-run/holos/pkg/wrapper"
)

//go:embed help/root.txt
var helpLong string

// Root holds the program configuration and root command for flag parsing.
type Root struct {
	Config  *app.Config
	Command *cobra.Command
}

// NewRoot builds a root cobra command with flags linked to the Config field.
func New(options ...app.Option) *cobra.Command {
	config := app.NewConfig(options...)
	root := Root{Config: config}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  helpLong,
		// We handle our own errors.
		SilenceUsage:  true,
		SilenceErrors: true,
		// Hidden because it annoys users trying to complete component
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			srvCfg := &server.Config{
				Port:          strconv.Itoa(config.ListenPort()),
				MetricsPort:   config.MetricsPort(),
				OIDCIssuer:    config.OIDCIssuer,
				OIDCAudiences: config.OIDCAudiences,
			}

			log := logger.FromContext(cmd.Context()).With("version", version.NewVersionInfo().Version)
			slog.SetDefault(log)
			log.DebugContext(ctx, "hello", "lifecycle", "start")

			app := app.App{
				Context: ctx,
				Logger:  log,
				Config:  config,
			}

			// Connect to the database
			conn, err := db.Client(app)
			if err != nil {
				return wrapper.Wrap(fmt.Errorf("could not create db client: %w", err))
			}
			defer func() {
				if closeError := conn.Client.Close(); closeError != nil {
					log.ErrorContext(ctx, "could not close database", "err", closeError)
				}
			}()

			// Retry until network is online or limit reached
			backoff := retry.NewFibonacci(1 * time.Second)
			backoff = retry.WithCappedDuration(5*time.Second, backoff)
			backoff = retry.WithMaxDuration(30*time.Second, backoff)
			// Ping the database
			ping := func(ctx context.Context) error {
				plog := slog.With("database", conn.Driver.Dialect(), "check", "network")
				plog.DebugContext(ctx, "ping")
				if pingErr := conn.DB.PingContext(ctx); pingErr != nil {
					plog.DebugContext(ctx, "retryable: could not ping", "ok", false, "err", pingErr)
					return retry.RetryableError(wrapper.Wrap(pingErr))
				}
				plog.DebugContext(ctx, "pong", "ok", true)
				return nil
			}
			if err = retry.Do(ctx, backoff, ping); err != nil {
				return wrapper.Wrap(err)
			}

			// Automatic migration
			if err = conn.Client.Schema.Create(ctx); err != nil {
				return wrapper.Wrap(err)
			}
			log.InfoContext(ctx, "schema created", "database", conn.Driver.Dialect())

			// Authentication (Identity Verifier)
			// We may pass an instrumented *http.Client via ctx in the future.
			verifier, err := authn.NewVerifier(app, config.OIDCIssuer)
			if err != nil {
				return wrapper.Wrap(fmt.Errorf("could not create identity verifier: %w", err))
			}

			// Start the server
			srv, err := server.NewServer(app, srvCfg, conn.Client, verifier)
			if err != nil {
				return wrapper.Wrap(fmt.Errorf("could not start server: %w", err))
			}

			if config.ListenAndServe {
				httpServer, healthy, ready := srv.ListenAndServe()
				stopCh := signals.SetupSignalHandler()
				sd := signals.NewShutdown(15*time.Second, log)
				sd.Graceful(stopCh, httpServer, healthy, ready)
			}

			return nil
		},
		// PersistentPreRunE runs after flag parsing before RunE.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Validate(cmd.OutOrStdout()); err != nil {
				return wrapper.Wrap(err)
			}

			level := config.GetLogLevel()
			var handler slog.Handler
			if config.LogFormat == "text" {
				handler = tint.NewHandler(cmd.ErrOrStderr(), &tint.Options{
					Level:       level,
					TimeFormat:  time.Kitchen,
					AddSource:   level == slog.LevelDebug,
					ReplaceAttr: config.ReplaceAttr,
				})
			} else {
				handler = slog.NewJSONHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{
					Level:       level,
					AddSource:   true,
					ReplaceAttr: config.ReplaceAttr,
				})
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Note, we deliberately do not touch slog.Default so these commands work concurrently.
			cmd.SetContext(logger.NewContext(ctx, slog.New(handler)))
			return nil
		},
	}

	// Add flags valid for all subcommands
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().AddGoFlagSet(config.FlagSet())
	root.Command = cmd
	return cmd
}
