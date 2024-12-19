package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/holos-run/holos/internal/cli"
	"github.com/holos-run/holos/internal/holos"
)

// MakeMain makes a main function for the cli or tests.
func MakeMain(options ...holos.Option) func() int {
	return func() (exitCode int) {
		cfg := holos.New(options...)
		slog.SetDefault(cfg.Logger())
		ctx := context.Background()

		if format := os.Getenv("HOLOS_CPU_PROFILE"); format != "" {
			f, _ := os.Create(fmt.Sprintf(format, os.Getppid(), os.Getpid()))
			err := pprof.StartCPUProfile(f)
			defer func() {
				pprof.StopCPUProfile()
				f.Close()
			}()
			if err != nil {
				return cli.HandleError(ctx, err, cfg)
			}
		}
		defer memProfile(ctx, cfg)

		if format := os.Getenv("HOLOS_TRACE"); format != "" {
			f, _ := os.Create(fmt.Sprintf(format, os.Getppid(), os.Getpid()))
			err := trace.Start(f)
			defer func() {
				trace.Stop()
				f.Close()
			}()
			if err != nil {
				return cli.HandleError(ctx, err, cfg)
			}
		}

		feature := &holos.EnvFlagger{}
		if err := cli.New(cfg, feature).ExecuteContext(ctx); err != nil {
			return cli.HandleError(ctx, err, cfg)
		}
		return 0
	}
}

func memProfile(ctx context.Context, cfg *holos.Config) {
	if format := os.Getenv("HOLOS_MEM_PROFILE"); format != "" {
		f, _ := os.Create(fmt.Sprintf(format, os.Getppid(), os.Getpid()))
		defer f.Close()
		if err := pprof.WriteHeapProfile(f); err != nil {
			_ = cli.HandleError(ctx, err, cfg)
		}
	}
}
