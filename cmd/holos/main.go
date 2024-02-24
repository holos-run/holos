package main

import (
	"context"
	"errors"
	"github.com/holos-run/holos/pkg/cli"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/holos-run/holos/pkg/wrapper"
	"log/slog"
	"os"
)

func main() {
	cfg := holos.New()
	slog.SetDefault(cfg.Logger())
	ctx := context.Background()
	if err := cli.New(cfg).ExecuteContext(ctx); err != nil {
		log := cfg.NewTopLevelLogger()
		var errAt *wrapper.ErrorAt
		const msg = "could not execute"
		if ok := errors.As(err, &errAt); ok {
			log.ErrorContext(ctx, msg, "err", errAt.Unwrap(), "loc", errAt.Source.Loc())
		} else {
			log.ErrorContext(ctx, msg, "err", err)
		}
		os.Exit(1)
	}
}
