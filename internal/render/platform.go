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
)

func Platform(ctx context.Context, pf *v1alpha1.Platform, stderr io.Writer) error {
	total := len(pf.Spec.Components)
	for idx, component := range pf.Spec.Components {
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
	}
	return nil
}
