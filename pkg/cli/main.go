package cli

import (
	"context"
	"errors"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/holos-run/holos/pkg/wrapper"
	"log/slog"
)

// MakeMain makes a main function for the cli or tests.
func MakeMain(options ...holos.Option) func() int {
	return func() (exitCode int) {
		cfg := holos.New(options...)
		slog.SetDefault(cfg.Logger())
		ctx := context.Background()
		if err := New(cfg).ExecuteContext(ctx); err != nil {
			return handleError(ctx, err, cfg)
		}
		return 0
	}
}

// handleError is the top level error handler that unwraps and logs errors.
func handleError(ctx context.Context, err error, hc *holos.Config) (exitCode int) {
	log := hc.NewTopLevelLogger()
	var errAt *wrapper.ErrorAt
	const msg = "could not execute"
	if ok := errors.As(err, &errAt); ok {
		log.ErrorContext(ctx, msg, "err", errAt.Unwrap(), "loc", errAt.Source.Loc())
	} else {
		log.ErrorContext(ctx, msg, "err", err)
	}
	return 1
}
