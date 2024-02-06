// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"context"
	"fmt"
	"github.com/holos-run/holos/pkg/wrapper"
)

type Builder struct {
	opts Options
}

// Options are options for a Builder.
// A zero Options consists entirely of default values.
type Options struct{}

// New returns a new *Builder with opts Options.
func New(opts Options) *Builder {
	return &Builder{opts: opts}
}

func (b *Builder) Run(ctx context.Context) error {
	return wrapper.Wrap(fmt.Errorf("not implemented"))
}
