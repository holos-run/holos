package cmd

import (
	"context"
	"log/slog"

	"github.com/holos-run/holos/internal/cli"
	"github.com/holos-run/holos/internal/holos"
)

// MakeMain makes a main function for the cli or tests.
func MakeMain(options ...holos.Option) func() int {
	return func() (exitCode int) {
		cfg := holos.New(options...)
		slog.SetDefault(cfg.Logger())
		ctx := context.Background()
		if err := cli.New(cfg).ExecuteContext(ctx); err != nil {
			return cli.HandleError(ctx, err, cfg)
		}
		return 0
	}
}
