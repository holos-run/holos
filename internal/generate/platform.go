package generate

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
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

func initPlatformMetadata(ctx context.Context, root, name string) error {
	log := logger.FromContext(ctx)
	rpcPlatform := &platform.Platform{Name: name}
	// Write the platform data.
	encoder := protojson.MarshalOptions{Indent: "  "}
	data, err := encoder.Marshal(rpcPlatform)
	if err != nil {
		return errors.Wrap(err)
	}
	if len(data) > 0 {
		data = append(data, '\n')
	}
	platformMetadataFile := filepath.Join(root, client.PlatformMetadataFile)
	if err := os.WriteFile(platformMetadataFile, data, 0o666); err != nil {
		return errors.Wrap(fmt.Errorf("could not write platform metadata: %w", err))
	}
	log.DebugContext(ctx, "wrote "+client.PlatformMetadataFile, "path", platformMetadataFile)

	return nil
}

// GeneratePlatform writes the cue code for a named platform to path.
func GeneratePlatform(ctx context.Context, dst, name string) error {
	log := logger.FromContext(ctx)
	// Check for a valid platform
	platformPath := filepath.Join(platformsRoot, name)
	if !dirExists(pfs, platformPath) {
		return errors.Wrap(fmt.Errorf("cannot generate: have: [%s] want: %+v", name, Platforms()))
	}

	platformMetadataFile := filepath.Join(dst, client.PlatformMetadataFile)
	if _, err := os.Stat(platformMetadataFile); err == nil {
		log.DebugContext(ctx, fmt.Sprintf("skipped write %s: already exists", platformMetadataFile))
	} else {
		if os.IsNotExist(err) {
			if err := initPlatformMetadata(ctx, dst, name); err != nil {
				return errors.Wrap(err)
			}
		} else {
			return errors.Wrap(err)
		}
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
