package generate

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
)

//go:embed all:platforms
var pfs embed.FS

// platformsRoot is the root path to copy platform cue code from.
const platformsRoot = "platforms"

// Platforms returns a slice of embedded platforms or nil if there are none.
func Platforms() []string {
	entries, err := fs.ReadDir(pfs, platformsRoot)
	if err != nil {
		return nil
	}
	dirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "cue.mod" {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs
}

// GeneratePlatform writes the cue code for a named platform to path.
func GeneratePlatform(ctx context.Context, dst, name string) error {
	log := logger.FromContext(ctx)
	// Check for a valid platform
	platformPath := filepath.Join(platformsRoot, name)
	if !dirExists(pfs, platformPath) {
		return errors.Wrap(fmt.Errorf("cannot generate: have: [%s] want: %+v", name, Platforms()))
	}

	// Copy the cue.mod directory
	if err := copyEmbedFS(ctx, pfs, filepath.Join(platformsRoot, "cue.mod"), filepath.Join(dst, "cue.mod"), bytes.NewBuffer); err != nil {
		return errors.Wrap(err)
	}

	// Copy the named platform
	if err := copyEmbedFS(ctx, pfs, platformPath, dst, bytes.NewBuffer); err != nil {
		return errors.Wrap(err)
	}

	log.DebugContext(ctx, "generated platform "+name, "path", getCwd(ctx))

	return nil
}
