// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"context"
	"fmt"
	"github.com/holos-run/holos/pkg/wrapper"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// A BuilderOption configures a Builder
type BuilderOption func(b *Builder)

type Builder struct {
	args []string
}

// Entrypoints are the leaf directories or files built by cue.
func Entrypoints(args []string) BuilderOption {
	return func(b *Builder) { b.args = args }
}

// New returns a new *Builder with opts Options.
func New(opts ...BuilderOption) *Builder {
	b := &Builder{}
	for _, option := range opts {
		option(b)
	}
	return b
}

func (b *Builder) Run(ctx context.Context) error {
	cueCtx := cuecontext.New()

	instances := load.Instances(b.args, nil)

	for _, instance := range instances {
		if err := instance.Err; err != nil {
			return wrapper.Wrap(fmt.Errorf("could not load: %w", err))
		}
		value := cueCtx.BuildInstance(instance)
		if err := value.Err(); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not build: %w", err))
		}
		if err := value.Validate(); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not validate: %w", err))
		}

		// Output in cue format
		fmt.Println(value)
	}

	return nil
}
