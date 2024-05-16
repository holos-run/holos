package generate

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
)

//go:embed all:platforms
var platforms embed.FS

// root is the root path to copy platform cue code from.
const root = "platforms"

// Platforms returns a slice of embedded platforms or nil if there are none.
func Platforms() []string {
	entries, err := fs.ReadDir(platforms, root)
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

// GeneratePlatform writes the cue code for a platform to the local working
// directory.
func GeneratePlatform(ctx context.Context, rpc *client.Client, orgID string, name string) error {
	log := logger.FromContext(ctx)
	// Check for a valid platform
	platformPath := filepath.Join(root, name)
	if !dirExists(platforms, platformPath) {
		return errors.Wrap(fmt.Errorf("cannot generate: have: [%s] want: %+v", name, Platforms()))
	}

	// Link the local platform the SaaS platform ID.
	rpcPlatforms, err := rpc.Platforms(ctx, orgID)
	if err != nil {
		return errors.Wrap(err)
	}

	var rpcPlatform *platform.Platform
	for _, p := range rpcPlatforms {
		if p.GetName() == name {
			rpcPlatform = p
			break
		}
	}
	if rpcPlatform == nil {
		return errors.Wrap(errors.New("cannot generate: platform not found in the holos server"))
	}

	// Write the platform data.
	data, err := json.MarshalIndent(rpcPlatform, "", "  ")
	if err != nil {
		return errors.Wrap(err)
	}
	if len(data) > 0 {
		data = append(data, '\n')
	}
	log = log.With("platform_id", rpcPlatform.GetId())
	path := "platform.metadata.json"
	if err := os.WriteFile(path, data, 0644); err != nil {
		return errors.Wrap(fmt.Errorf("could not write platform metadata: %w", err))
	}
	log.InfoContext(ctx, "wrote "+path, "path", filepath.Join(getCwd(ctx), path))

	// Copy the cue.mod directory
	if err := copyEmbedFS(ctx, platforms, filepath.Join(root, "cue.mod"), "cue.mod"); err != nil {
		return errors.Wrap(err)
	}

	// Copy the named platform
	if err := copyEmbedFS(ctx, platforms, platformPath, "."); err != nil {
		return errors.Wrap(err)
	}

	log.InfoContext(ctx, "generated platform "+name, "path", getCwd(ctx))

	return nil
}

func dirExists(srcFS embed.FS, path string) bool {
	entries, err := fs.ReadDir(srcFS, path)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

func copyEmbedFS(ctx context.Context, srcFS embed.FS, srcPath, dstPath string) error {
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
		} else {
			data, err := srcFS.ReadFile(path)
			if err != nil {
				return errors.Wrap(err)
			}
			if err := os.WriteFile(dstFullPath, data, os.ModePerm); err != nil {
				return errors.Wrap(err)
			}
			log.DebugContext(ctx, "wrote", "file", dstFullPath)
		}
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
