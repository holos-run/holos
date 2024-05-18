package client

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogo/protobuf/jsonpb"
	platform "github.com/holos-run/holos/service/gen/holos/platform/v1alpha1"
)

// LoadPlatform loads the platform.metadata.json file from a named path.  Useful
// to obtain a platform id for PlatformService rpc methods.
func LoadPlatform(ctx context.Context, name string) (*platform.Platform, error) {
	data, err := os.ReadFile(filepath.Join(name, "platform.metadata.json"))
	if err != nil {
		return nil, fmt.Errorf("could not load platform metadata: %w", err)
	}
	p := &platform.Platform{}
	if err := jsonpb.Unmarshal(bytes.NewReader(data), p); err != nil {
		return nil, fmt.Errorf("could not load platform metadata: %w", err)
	}
	return p, nil
}
