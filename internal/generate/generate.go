package generate

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
)

func dirExists(srcFS embed.FS, path string) bool {
	entries, err := fs.ReadDir(srcFS, path)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// copyEmbedFS copies embedded files from srcPath to dstPath passing the
// contents through mapFunc.
func copyEmbedFS(ctx context.Context, srcFS embed.FS, srcPath, dstPath string, mapFunc func([]byte) *bytes.Buffer) error {
	log := logger.FromContext(ctx)
	return fs.WalkDir(srcFS, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err)
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return errors.Wrap(err)
		}

		dstFullPath := filepath.Join(dstPath, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(dstFullPath, os.ModePerm); err != nil {
				return errors.Wrap(err)
			}
			log.DebugContext(ctx, "created", "directory", dstFullPath)
			return nil
		}

		if filepath.Base(path) == "schematic.json" {
			log.DebugContext(ctx, "skipped", "file", dstFullPath)
			return nil
		}

		data, err := srcFS.ReadFile(path)
		if err != nil {
			return errors.Wrap(err)
		}

		buf := mapFunc(data)
		if err := os.WriteFile(dstFullPath, buf.Bytes(), 0o666); err != nil {
			return errors.Wrap(err)
		}

		log.DebugContext(ctx, "wrote", "file", dstFullPath)
		return nil
	})
}

func getCwd(ctx context.Context) string {
	cwd, err := os.Getwd()
	if err != nil {
		logger.FromContext(ctx).WarnContext(ctx, "could not get working directory", "err", err)
		return "."
	}
	abs, err := filepath.Abs(cwd)
	if err != nil {
		logger.FromContext(ctx).WarnContext(ctx, "could not get absolute path", "err", err)
		return cwd
	}
	return abs
}
