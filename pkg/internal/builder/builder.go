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

// An Option configures a Builder
type Option func(*config)

type config struct {
	args []string
}

type Builder struct {
	cfg config
}

// New returns a new *Builder configured by opts Option.
func New(opts ...Option) *Builder {
	var cfg config
	for _, f := range opts {
		f(&cfg)
	}
	b := &Builder{cfg: cfg}
	return b
}

// Entrypoints configures the leaf directories Builder builds.
func Entrypoints(args []string) Option {
	return func(cfg *config) { cfg.args = args }
}

type buildInfo struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

type out struct {
	Out string `json:"out,omitempty"`
}

func (b *Builder) Run(ctx context.Context) error {
	cueCtx := cuecontext.New()

	instances := load.Instances(b.cfg.args, nil)

	for _, instance := range instances {
		var info buildInfo
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

		if err := value.Decode(&info); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
		}

		switch kind := info.Kind; kind {
		case "KubernetesObjects":
			var out out
			if err := value.Decode(&out); err != nil {
				return wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
			}
			fmt.Printf(out.Out)
		default:
			return wrapper.Wrap(fmt.Errorf("build kind not implemented: %v", kind))
		}
	}

	return nil
}
