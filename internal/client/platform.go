package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/server/middleware/logger"
	object "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
)

// PlatformMetadataFile is the platform metadata json file name located in the root
// of a platform directory.
const PlatformMetadataFile = "platform.metadata.json"

// PlatformConfigFile is the marshaled json representation of the PlatformConfig
// DTO used to cache the data holos passes from the PlatformService to CUE when
// rendering platform components.
const PlatformConfigFile = "platform.config.json"

// LoadPlatform loads the platform.metadata.json file from a named path.  Useful
// to obtain a platform id for PlatformService rpc methods.
func LoadPlatform(ctx context.Context, name string) (*platform.Platform, error) {
	data, err := os.ReadFile(filepath.Join(name, PlatformMetadataFile))
	if err != nil {
		return nil, fmt.Errorf("could not load platform metadata: %w", err)
	}
	p := &platform.Platform{}
	if err := protojson.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("could not load platform metadata: %w", err)
	}
	return p, nil
}

// LoadPlatformConfig loads the PlatformConfig DTO from the platform.config.json
// file.  Useful to provide all values necessary to render cue config without an
// rpc to the HolosService.
func LoadPlatformConfig(ctx context.Context, name string) (*object.PlatformConfig, error) {
	data, err := os.ReadFile(filepath.Join(name, PlatformConfigFile))
	if err != nil {
		return nil, fmt.Errorf("could not load platform config: %w", err)
	}
	pc := &object.PlatformConfig{}
	if err := protojson.Unmarshal(data, pc); err != nil {
		return nil, fmt.Errorf("could not load platform config: %w", err)
	}
	return pc, nil
}

// SavePlatformConfig writes pc to the platform root directory path identified by name.
func SavePlatformConfig(ctx context.Context, name string, pc *object.PlatformConfig) (string, error) {
	data, err := protojson.Marshal(pc)
	if err != nil {
		return "", err
	}
	if len(data) > 0 {
		data = append(data, '\n')
	}
	path := filepath.Join(name, PlatformConfigFile)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("could not write platform config: %w", err)
	}
	logger.FromContext(ctx).DebugContext(ctx, "wrote", "path", path)
	return path, nil
}
