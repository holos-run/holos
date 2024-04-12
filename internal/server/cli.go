package server

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/spf13/cobra"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/db"
	"github.com/holos-run/holos/internal/server/middleware/authn"
	"github.com/holos-run/holos/internal/server/server"
	"github.com/holos-run/holos/internal/server/signals"
)

//go:embed help/root.txt
var helpLong string

// New builds a root cobra command with flags linked to the Config field.
func New(cfg *holos.Config) *cobra.Command {
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

			log := cfg.Logger()
			log.DebugContext(ctx, "hello", "lifecycle", "start")

			// Connect to the database
			conn, err := db.Client(cfg)
			if err != nil {
				return errors.Wrap(fmt.Errorf("could not create db client: %w", err))
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
					return retry.RetryableError(errors.Wrap(pingErr))
				}
				plog.DebugContext(ctx, "pong", "ok", true)
				return nil
			}
			if err = retry.Do(ctx, backoff, ping); err != nil {
				return errors.Wrap(err)
			}

			// Automatic migration
			if err = conn.Client.Schema.Create(ctx); err != nil {
				return errors.Wrap(err)
			}
			log.InfoContext(ctx, "schema created", "database", conn.Driver.Dialect())

			// Authentication (Identity Verifier)
			// We may pass an instrumented *http.Client via ctx in the future.
			verifier, err := authn.NewVerifier(ctx, log, cfg.ServerConfig.OIDCIssuer())
			if err != nil {
				return errors.Wrap(fmt.Errorf("could not create identity verifier: %w", err))
			}

			// Start the server
			srv, err := server.NewServer(cfg, conn.Client, verifier)
			if err != nil {
				return errors.Wrap(fmt.Errorf("could not start server: %w", err))
			}

			if cfg.ServerConfig.ListenAndServe() {
				httpServer, healthy, ready := srv.ListenAndServe()
				stopCh := signals.SetupSignalHandler()
				sd := signals.NewShutdown(15*time.Second, log)
				sd.Graceful(stopCh, httpServer, healthy, ready)
			}

			return nil
		},
	}

	// Add flags valid for all subcommands
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().AddGoFlagSet(cfg.ServerFlagSet())
	return cmd
}
