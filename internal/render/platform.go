package render

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/holos-run/holos/api/v1alpha1"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/util"
	"golang.org/x/sync/errgroup"
)

func Platform(ctx context.Context, concurrency int, pf *v1alpha1.Platform, stderr io.Writer) error {
	total := len(pf.Spec.Components)

	g, ctx := errgroup.WithContext(ctx)
	// Limit the number of concurrent goroutines due to CUE memory usage concerns
	// while rendering components.  One more for the producer.
	g.SetLimit(concurrency + 1)

	// Spawn a producer because g.Go() blocks when the group limit is reached.
	g.Go(func() error {
		for idx, component := range pf.Spec.Components {
			select {
			case <-ctx.Done():
				return nil
			default:
				// Capture idx and component to avoid issues with closure. Can be removed on Go 1.22.
				idx, component := idx, component
				// Worker go routine.  Blocks if limit has been reached.
				g.Go(func() error {
					start := time.Now()
					log := logger.FromContext(ctx).With("path", component.Path, "cluster", component.Cluster, "num", idx+1, "total", total)
					log.DebugContext(ctx, "render component")

					// Execute a sub-process to limit CUE memory usage.
					args := []string{"render", "component", "--cluster-name", component.Cluster, component.Path}
					result, err := util.RunCmd(ctx, "holos", args...)
					if err != nil {
						_, _ = io.Copy(stderr, result.Stderr)
						return errors.Wrap(fmt.Errorf("could not render component: %w", err))
					}

					duration := time.Since(start)
					log.InfoContext(ctx, "ok render component", "duration", duration)
					return nil
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	return g.Wait()
}
