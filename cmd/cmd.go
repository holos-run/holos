package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/holos-run/holos/internal/cli"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/version"
)

// MakeMain makes a main function for the cli or tests.
func MakeMain(options ...holos.Option) func() int {
	return func() (exitCode int) {
		// TODO(jjm): check HOLOS_CHDIR and chdir if set for tests.
		if len(os.Args) >= 2 && os.Args[1] == "version" {
			if _, err := fmt.Println(version.GetVersion()); err != nil {
				panic(err)
			}
			return 0
		}
		cfg := holos.New(options...)
		slog.SetDefault(cfg.Logger())
		ctx := context.Background()
		cmd := cli.New(cfg)

		if err := cmd.ExecuteContext(ctx); err != nil {
			return cli.HandleError(ctx, err, cfg)
		}
		return 0
	}
}
