package v1alpha4

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/holos-run/holos/api/core/v1alpha4"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/render"
	"github.com/holos-run/holos/internal/util"
	"golang.org/x/sync/errgroup"
)

// Platform represents a platform builder.
type PlatformBuilder struct {
	Platform    v1alpha4.Platform
	Concurrency int
	Stderr      io.Writer
}

func (b *PlatformBuilder) Build(ctx context.Context, s render.Setter) error {
	parentStart := time.Now()
	components := b.Platform.Spec.Components
	total := len(components)
	g, ctx := errgroup.WithContext(ctx)
	// Limit the number of concurrent goroutines due to CUE memory usage concerns
	// while rendering components.  One more for the producer.
	g.SetLimit(b.Concurrency + 1)
	// Spawn a producer because g.Go() blocks when the group limit is reached.
	g.Go(func() error {
		for idx, component := range components {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Capture idx and component to avoid issues with closure. Can be
				// removed on Go 1.22.
				idx, component := idx, component
				// Worker go routine.  Blocks if limit has been reached.
				g.Go(func() error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						start := time.Now()
						log := logger.FromContext(ctx).With(
							"path", component.Path,
							"cluster", component.Cluster,
							"num", idx+1,
							"total", total,
						)
						log.DebugContext(ctx, "render component")

						// Execute a sub-process to limit CUE memory usage.
						args := []string{
							"render",
							"component",
							"--cluster-name",
							component.Cluster,
							component.Path,
						}
						result, err := util.RunCmd(ctx, "holos", args...)
						if err != nil {
							_, _ = io.Copy(b.Stderr, result.Stderr)
							return errors.Wrap(fmt.Errorf("could not render component: %w", err))
						}

						duration := time.Since(start)
						msg := fmt.Sprintf(
							"rendered %s for cluster %s in %s",
							filepath.Base(component.Path),
							component.Cluster,
							duration,
						)
						log.InfoContext(ctx, msg, "duration", duration)
						return nil
					}
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Since(parentStart)
	msg := fmt.Sprintf("rendered platform in %s", duration)
	logger.FromContext(ctx).InfoContext(ctx, msg, "duration", duration, "version", b.Platform.APIVersion)
	return nil
}
