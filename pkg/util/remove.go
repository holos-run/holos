package util

import (
	"context"
	"github.com/holos-run/holos/pkg/logger"
	"os"
)

// Remove cleans up path, useful for temporary directories.
func Remove(ctx context.Context, path string) {
	log := logger.FromContext(ctx)
	if err := os.RemoveAll(path); err != nil {
		log.WarnContext(ctx, "tmp: could not Remove", "err", err, "path", path)
	} else {
		log.DebugContext(ctx, "tmp: removed", "path", path)
	}
}
