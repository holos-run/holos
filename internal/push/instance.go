package push

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/errors"
)

func NewInstance(name string) (*Instance, error) {
	absPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	path := holos.InstancePath(absPath)
	mod, err := FindCueMod(path)
	if err != nil {
		return nil, err
	}
	return &Instance{path: path, mod: mod}, nil
}

// Instance represents a CUE instance.
type Instance struct {
	path holos.InstancePath
	mod  holos.PathCueMod
}

// Export builds the cue instance into a JSON byte slice.  Equivalent of cue
// export.
func (i *Instance) Export(ctx context.Context) ([]byte, error) {
	// CUE Loader
	cfg := load.Config{Dir: string(i.mod)}

	// Make target relative to the module directory
	relPath, err := filepath.Rel(string(i.mod), string(i.path))
	if err != nil {
		return nil, err
	}
	relPath = "./" + relPath

	instances := load.Instances([]string{relPath}, &cfg)
	if len(instances) != 1 {
		return nil, errors.Wrap(errors.New("exactly one instance is required"))
	}
	instance := instances[0]

	if err := instance.Err; err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not load: %w", err))
	}

	cueCtx := cuecontext.New()

	value := cueCtx.BuildInstance(instance)
	if err := value.Err(); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not build %s: %w", instance.Dir, err))
	}
	if err := value.Validate(); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not validate: %w", err))
	}

	b, err := value.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not marshal cue instance %s: %w", instance.Dir, err))
	}

	return b, nil
}

func FindCueMod(name holos.InstancePath) (dir holos.PathCueMod, err error) {
	path := holos.PathCueMod(name)
	for {
		if _, err := os.Stat(filepath.Join(string(path), "cue.mod")); err == nil {
			dir = path
			break
		} else if !os.IsNotExist(err) {
			return "", err
		}
		parentPath := holos.PathCueMod(filepath.Dir(string(path)))
		if parentPath == path {
			return "", fmt.Errorf("no cue.mod from root to leaf: %v", name)
		}
		path = parentPath
	}
	return dir, nil
}
