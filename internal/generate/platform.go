package generate

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
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
	return generatePlatform(ctx, dst, name, false)
}

// GeneratePlatformWithSDK writes a platform together with the legacy CUE SDK.
// It is reserved for internal compatibility fixtures that intentionally import
// the pre-v1beta1 author packages. Fresh v1beta1 platforms use the compact
// local runtime from GeneratePlatform instead.
func GeneratePlatformWithSDK(ctx context.Context, dst, name string) error {
	return generatePlatform(ctx, dst, name, true)
}

func generatePlatform(ctx context.Context, dst, name string, withSDK bool) error {
	log := logger.FromContext(ctx)
	// Check for a valid platform
	platformPath := filepath.Join(platformsRoot, name)
	if !dirExists(pfs, platformPath) {
		return errors.Wrap(fmt.Errorf("cannot generate: have: [%s] want: %+v", name, Platforms()))
	}

	if name == "v1beta1" && !withSDK {
		// The v1beta1 generated platform has a deliberately small local runtime.
		// Copy only its module declaration; copying cue.mod/gen and cue.mod/pkg
		// would vendor the complete Holos CUE SDK into every initialized platform.
		modulePath := filepath.Join(platformsRoot, "cue.mod", "module.cue")
		module, err := pfs.ReadFile(modulePath)
		if err != nil {
			return errors.Wrap(err)
		}
		if err := os.MkdirAll(filepath.Join(dst, "cue.mod"), os.ModePerm); err != nil {
			return errors.Wrap(err)
		}
		if err := os.WriteFile(filepath.Join(dst, "cue.mod", "module.cue"), module, 0o666); err != nil {
			return errors.Wrap(err)
		}
	} else if err := copyEmbedFS(ctx, pfs, filepath.Join(platformsRoot, "cue.mod"), filepath.Join(dst, "cue.mod"), bytes.NewBuffer); err != nil {
		return errors.Wrap(err)
	}

	// Copy the named platform
	if err := copyEmbedFS(ctx, pfs, platformPath, dst, bytes.NewBuffer); err != nil {
		return errors.Wrap(err)
	}

	log.DebugContext(ctx, "generated platform "+name, "path", getCwd(ctx))

	return nil
}
