package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/holos-run/holos"
	core "github.com/holos-run/holos/api/core/v1alpha2"
	meta "github.com/holos-run/holos/api/meta/v1alpha2"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
)

// Platform builds a platform
func (b *Builder) Platform(ctx context.Context, cfg *client.Config) (*core.Platform, error) {
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "cue: building platform instance")
	bd, err := b.Unify(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return b.runPlatform(ctx, bd)
}

func (b Builder) runPlatform(ctx context.Context, bd BuildData) (*core.Platform, error) {
	path := holos.InstancePath(bd.Dir)
	log := logger.FromContext(ctx).With("dir", path)

	value := bd.Value
	if err := bd.Value.Err(); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not load: %w", err))
	}

	log.DebugContext(ctx, "cue: validating instance")
	if err := value.Validate(); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not validate: %w", err))
	}

	log.DebugContext(ctx, "cue: decoding holos platform")
	jsonBytes, err := value.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not marshal cue instance %s: %w", bd.Dir, err))
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	// Discriminate the type of build plan.
	tm := &meta.TypeMeta{}
	err = decoder.Decode(tm)
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("invalid platform: %s: %w", bd.Dir, err))
	}

	log.DebugContext(ctx, "cue: discriminated build kind: "+tm.GetKind(), "kind", tm.GetKind(), "apiVersion", tm.GetAPIVersion())

	// New decoder for the full object
	decoder = json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.DisallowUnknownFields()

	var pf core.Platform
	switch tm.GetKind() {
	case "Platform":
		if err = decoder.Decode(&pf); err != nil {
			err = errors.Wrap(fmt.Errorf("could not decode platform %s: %w", bd.Dir, err))
			return nil, err
		}
		return &pf, nil
	default:
		err = errors.Wrap(fmt.Errorf("unknown kind: %v", tm.GetKind()))
	}

	return nil, err
}
