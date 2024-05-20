package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"github.com/holos-run/holos"
	"github.com/holos-run/holos/api/v1alpha1"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
)

// Platform builds a platform
func (b *Builder) Platform(ctx context.Context, cfg *client.Config) (*v1alpha1.Platform, error) {
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "cue: building platform instance")
	instances, err := b.Instances(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// We only process the first instance, assume the render platform subcommand enforces this.
	for idx, instance := range instances {
		log.DebugContext(ctx, "cue: building instance", "idx", idx, "dir", instance.Dir)
		p, err := b.runPlatform(ctx, instance)
		if err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not build platform: %w", err))
		}
		return p, nil
	}

	return nil, errors.Wrap(errors.New("missing platform instance"))
}

func (b Builder) runPlatform(ctx context.Context, instance *build.Instance) (*v1alpha1.Platform, error) {
	path := holos.InstancePath(instance.Dir)
	log := logger.FromContext(ctx).With("dir", path)

	if err := instance.Err; err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not load: %w", err))
	}
	cueCtx := cuecontext.New()
	value := cueCtx.BuildInstance(instance)
	if err := value.Err(); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not build %s: %w", instance.Dir, err))
	}
	log.DebugContext(ctx, "cue: validating instance")
	if err := value.Validate(); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not validate: %w", err))
	}

	log.DebugContext(ctx, "cue: decoding holos platform")
	jsonBytes, err := value.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not marshal cue instance %s: %w", instance.Dir, err))
	}
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	// Discriminate the type of build plan.
	tm := &v1alpha1.TypeMeta{}
	err = decoder.Decode(tm)
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("invalid platform: %s: %w", instance.Dir, err))
	}

	log.DebugContext(ctx, "cue: discriminated build kind: "+tm.Kind, "kind", tm.Kind, "apiVersion", tm.APIVersion)

	// New decoder for the full object
	decoder = json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.DisallowUnknownFields()

	var pf v1alpha1.Platform
	switch tm.Kind {
	case "Platform":
		if err = decoder.Decode(&pf); err != nil {
			err = errors.Wrap(fmt.Errorf("could not decode platform %s: %w", instance.Dir, err))
			return nil, err
		}
		return &pf, nil
	default:
		err = errors.Wrap(fmt.Errorf("unknown kind: %v", tm.Kind))
	}

	return nil, err
}
