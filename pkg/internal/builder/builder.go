// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"context"
	"cuelang.org/go/cue/load"
	"fmt"
	"github.com/holos-run/holos/pkg/wrapper"

	"cuelang.org/go/cue/cuecontext"
	// "cuelang.org/go/cue/load"
)

type Builder struct {
	opts Options
}

// Options are options for a Builder.
// A zero Options consists entirely of default values.
type Options struct {
	// Entrypoints are the cue entrypoints, same as are passed to the cue cli
	Entrypoints []string
}

// New returns a new *Builder with opts Options.
func New(opts Options) *Builder {
	return &Builder{opts: opts}
}

func (b *Builder) Run(ctx context.Context) error {
	cueCtx := cuecontext.New()

	buildInstances := load.Instances(b.opts.Entrypoints, nil)

	for _, bi := range buildInstances {
		if err := bi.Err; err != nil {
			return wrapper.Wrap(fmt.Errorf("could not load: %w", err))
		}
		value := cueCtx.BuildInstance(bi)
		if err := value.Err(); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not build: %w", err))
		}
		if err := value.Validate(); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not validate: %w", err))
		}

		fmt.Println("value:", value)
	}

	return nil
}
