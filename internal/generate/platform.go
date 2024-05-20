package generate

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
)

//go:embed all:platforms
var platforms embed.FS

// platformsRoot is the root path to copy platform cue code from.
const platformsRoot = "platforms"

// Platforms returns a slice of embedded platforms or nil if there are none.
func Platforms() []string {
	entries, err := fs.ReadDir(platforms, platformsRoot)
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
	platformPath := filepath.Join(platformsRoot, name)
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
	if err := os.WriteFile(client.PlatformMetadataFile, data, 0644); err != nil {
		return errors.Wrap(fmt.Errorf("could not write platform metadata: %w", err))
	}
	log.InfoContext(ctx, "wrote "+client.PlatformMetadataFile, "path", filepath.Join(getCwd(ctx), client.PlatformMetadataFile))

	// Copy the cue.mod directory
	if err := copyEmbedFS(ctx, platforms, filepath.Join(platformsRoot, "cue.mod"), "cue.mod", bytes.NewBuffer); err != nil {
		return errors.Wrap(err)
	}

	// Copy the named platform
	if err := copyEmbedFS(ctx, platforms, platformPath, ".", bytes.NewBuffer); err != nil {
		return errors.Wrap(err)
	}

	log.InfoContext(ctx, "generated platform "+name, "path", getCwd(ctx))

	return nil
}
