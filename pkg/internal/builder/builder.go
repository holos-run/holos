// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"context"
	"fmt"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/wrapper"
	"os"
	"path/filepath"

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
	log := logger.FromContext(ctx)
	cueCtx := cuecontext.New()

	dir, err := b.findCueMod()
	if err != nil {
		return wrapper.Wrap(err)
	}

	cfg := load.Config{Dir: dir}

	// Make args relative to the module directory
	args := make([]string, len(b.cfg.args))
	for idx, path := range b.cfg.args {
		target, err := filepath.Abs(path)
		if err != nil {
			return wrapper.Wrap(fmt.Errorf("could not find absolute path: %w", err))
		}
		relPath, err := filepath.Rel(dir, target)
		if err != nil {
			return wrapper.Wrap(fmt.Errorf("invalid argument, must be relative to cue.mod: %w", err))
		}
		relPath = "./" + relPath
		args[idx] = relPath
		equiv := fmt.Sprintf("(cd %v && cue export --out text -e out %v)", dir, relPath)
		log.Debug("equivalent command", "cue", equiv)
	}

	instances := load.Instances(args, &cfg)

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

// findCueMod returns the root module location containing the cue.mod file or
// directory or an error if the builder arguments do not share a common root
// module.
func (b *Builder) findCueMod() (dir string, err error) {
	for _, origPath := range b.cfg.args {
		var path string
		if path, err = filepath.Abs(origPath); err != nil {
			return
		}
		for {
			if _, err := os.Stat(filepath.Join(path, "cue.mod")); err == nil {
				if dir != "" && dir != path {
					return "", fmt.Errorf("multiple modules not supported: %v is not %v", dir, path)
				}
				dir = path
				break
			} else if !os.IsNotExist(err) {
				return "", err
			}
			parentPath := filepath.Dir(path)
			if parentPath == path {
				return "", fmt.Errorf("no cue.mod from root to leaf: %v", origPath)
			}
			path = parentPath
		}
	}
	return dir, nil
}
