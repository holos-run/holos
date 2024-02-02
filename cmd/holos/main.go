package main

import (
	"context"
	"github.com/holos-run/holos/pkg/cli"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/wrapper"
	"log/slog"
	"os"
)

func main() {
	cfg := config.New(os.Stdout, os.Stderr)
	slog.SetDefault(cfg.Logger())
	ctx := context.Background()
	if err := cli.New(cfg).ExecuteContext(ctx); err != nil {
		wrapper.LogError(ctx, err)
		os.Exit(1)
	}
}
