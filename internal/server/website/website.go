// Package website serves the holos website with a http server.
package website

import (
	"context"
	_ "embed"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/holos-run/holos/doc/website"
	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
)

//go:embed help.txt
var helpLong string

type config struct {
	addr string
}

func (c *config) Addr() string {
	if c == nil {
		return ""
	}
	return c.addr
}

// New builds a root cobra command with flags linked to the Config field.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("website")
	cmd.Short = "serve the holos website"
	cmd.Long = helpLong
	cmd.Args = cobra.NoArgs

	wcfg := &config{}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&wcfg.addr, "http", ":3001", "http service address")

	cmd.Flags().AddGoFlagSet(fs)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context(), wcfg)
	}

	return cmd
}

func run(ctx context.Context, cfg Config) error {
	if ctx == nil {
		return errors.Wrap(errors.New("could not run: no context"))
	}

	log := logger.FromContext(ctx).With("addr", cfg.Addr())

	mux := http.NewServeMux()
	fs := http.FS(website.Root())
	mux.Handle("/", http.FileServer(fs))

	// Wrapped mux with middleware
	wmux := logger.ResponseLogger(logger.LoggingMiddleware(log)(mux))

	server := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           h2c.NewHandler(wmux, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	var g errgroup.Group
	g.Go(func() error {
		return server.ListenAndServe()
	})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		<-signalChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})

	log.InfoContext(ctx, "listening for http requests")

	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

type Config interface {
	Addr() string
}
