package cli

import (
	"context"
	"fmt"
	"log/slog"

	cue "cuelang.org/go/cue/errors"
	"github.com/holos-run/holos/pkg/errors"
	"github.com/holos-run/holos/pkg/holos"
)

// MakeMain makes a main function for the cli or tests.
func MakeMain(options ...holos.Option) func() int {
	return func() (exitCode int) {
		cfg := holos.New(options...)
		slog.SetDefault(cfg.Logger())
		ctx := context.Background()
		if err := New(cfg).ExecuteContext(ctx); err != nil {
			return HandleError(ctx, err, cfg)
		}
		return 0
	}
}

// HandleError is the top level error handler that unwraps and logs errors.
func HandleError(ctx context.Context, err error, hc *holos.Config) (exitCode int) {
	log := hc.NewTopLevelLogger()
	var cueErr cue.Error
	var errAt *errors.ErrorAt
	const msg = "could not execute"
	if errors.As(err, &errAt) {
		log.ErrorContext(ctx, msg, "err", errAt.Unwrap(), "loc", errAt.Source.Loc())
	} else {
		log.ErrorContext(ctx, msg, "err", err)
	}
	// cue errors are bundled up as a list and refer to multiple files / lines.
	if errors.As(err, &cueErr) {
		msg := cue.Details(cueErr, nil)
		_, _ = fmt.Fprint(hc.Stderr(), msg)
	}
	return 1
}
