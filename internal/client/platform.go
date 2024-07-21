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

// PlatformMetadataFile is the platform metadata json file name located in the
// root of a platform directory.  This file is the authoritative source of truth
// for the PlatformID used in rpc calls to the PlatformService.
const PlatformMetadataFile = "platform.metadata.json"

// PlatformConfigFile represents the marshaled json representation of the
// PlatformConfig DTO used to persist the inputs to the CUE platform code.
const PlatformConfigFile = "platform.config.json"

// LoadPlatformMetadata loads the platform.metadata.json file from a named path.
// Used as the authoritative source of truth to obtain a platform id for
// PlatformService rpc methods.
func LoadPlatformMetadata(ctx context.Context, name string) (*platform.Platform, error) {
	data, err := os.ReadFile(filepath.Join(name, PlatformMetadataFile))
	if err != nil {
		return nil, fmt.Errorf("could not load platform metadata: %w", err)
	}
	p := &platform.Platform{}
	if err := protojson.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("could not load platform metadata: %w", err)
	}
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "loaded: "+p.GetName(), "path", PlatformMetadataFile, "name", p.GetName(), "display_name", p.GetDisplayName(), "id", p.GetId())
	return p, nil
}

// LoadPlatformConfig loads the PlatformConfig DTO from the platform.config.json
// file.  Useful to provide all values necessary to render cue config without an
// rpc to the HolosService.
func LoadPlatformConfig(ctx context.Context, name string) (*object.PlatformConfig, error) {
	data, err := os.ReadFile(filepath.Join(name, PlatformConfigFile))
	if err != nil {
		return nil, fmt.Errorf("could not load platform model: %w", err)
	}
	p := &object.PlatformConfig{}
	if err := protojson.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("could not load platform model: %w", err)
	}
	return p, nil
}

// SavePlatformConfig writes pc to the platform root directory path identified by name.
func SavePlatformConfig(ctx context.Context, name string, pc *object.PlatformConfig) (string, error) {
	encoder := protojson.MarshalOptions{Multiline: true}
	data, err := encoder.Marshal(pc)
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
