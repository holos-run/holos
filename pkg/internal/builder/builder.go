// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/holos-run/holos/api/v1alpha1"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/wrapper"
)

const (
	KubernetesObjects = v1alpha1.KubernetesObjectsKind
	// Helm is the value of the kind field of holos build output indicating helm
	// values and helm command information.
	Helm = v1alpha1.HelmChartKind
	// Skip is the value when the instance should be skipped
	Skip = "Skip"
	// KustomizeBuild is the value of the kind field of cue output indicating holos should process the component using kustomize build to render output.
	KustomizeBuild = v1alpha1.KustomizeBuildKind
)

// An Option configures a Builder
type Option func(*config)

type config struct {
	args    []string
	cluster string
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

// Cluster configures the cluster name for the holos component instance.
func Cluster(name string) Option {
	return func(cfg *config) { cfg.cluster = name }
}

// Cluster returns the cluster name of the component instance being built.
func (b *Builder) Cluster() string {
	return b.cfg.cluster
}

// Instances returns the cue build instances being built.
func (b *Builder) Instances(ctx context.Context) ([]*build.Instance, error) {
	log := logger.FromContext(ctx)

	mod, err := b.findCueMod()
	if err != nil {
		return nil, wrapper.Wrap(err)
	}
	dir := string(mod)

	cfg := load.Config{Dir: dir}

	// Make args relative to the module directory
	args := make([]string, len(b.cfg.args))
	for idx, path := range b.cfg.args {
		target, err := filepath.Abs(path)
		if err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not find absolute path: %w", err))
		}
		relPath, err := filepath.Rel(dir, target)
		if err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("invalid argument, must be relative to cue.mod: %w", err))
		}
		relPath = "./" + relPath
		args[idx] = relPath
		equiv := fmt.Sprintf("cue export --out yaml -t cluster=%v %v", b.Cluster(), relPath)
		log.Debug("cue: equivalent command: " + equiv)
	}

	// Refer to https://github.com/cue-lang/cue/blob/v0.7.0/cmd/cue/cmd/common.go#L429
	cfg.Tags = append(cfg.Tags, "cluster="+b.Cluster())
	log.DebugContext(ctx, fmt.Sprintf("cue: tags %v", cfg.Tags))

	return load.Instances(args, &cfg), nil
}

func (b *Builder) Run(ctx context.Context) (results []*v1alpha1.Result, err error) {
	results = make([]*v1alpha1.Result, 0, len(b.cfg.args))
	cueCtx := cuecontext.New()
	logger.FromContext(ctx).DebugContext(ctx, "cue: building instances")
	instances, err := b.Instances(ctx)
	if err != nil {
		return results, err
	}

	// Each CUE instance provides a BuildPlan
	for _, instance := range instances {
		var buildPlan v1alpha1.BuildPlan

		log := logger.FromContext(ctx).With("dir", instance.Dir)
		if err := instance.Err; err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not load: %w", err))
		}
		log.DebugContext(ctx, "cue: building instance")
		value := cueCtx.BuildInstance(instance)
		if err := value.Err(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not build %s: %w", instance.Dir, err))
		}
		log.DebugContext(ctx, "cue: validating instance")
		if err := value.Validate(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not validate: %w", err))
		}

		log.DebugContext(ctx, "cue: decoding holos build plan")
		// Hack to catch unknown fields https://github.com/holos-run/holos/issues/72
		jsonBytes, err := value.MarshalJSON()
		if err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not marshal cue instance %s: %w", instance.Dir, err))
		}
		decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&buildPlan)
		if err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("invalid BuildPlan: %s: %w", instance.Dir, err))
		}

		if err := buildPlan.Validate(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not validate %s: %w", instance.Dir, err))
		}

		if buildPlan.Spec.Disabled {
			log.DebugContext(ctx, "skipped: spec.disabled is true", "skipped", true)
			continue
		}

		// TODO: concurrent renders
		for _, component := range buildPlan.Spec.Components.KubernetesObjectsList {
			if result, err := component.Render(ctx, holos.PathComponent(instance.Dir)); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not render: %w", err))
			} else {
				results = append(results, result)
			}
		}
		for _, component := range buildPlan.Spec.Components.HelmChartList {
			if result, err := component.Render(ctx, holos.PathComponent(instance.Dir)); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not render: %w", err))
			} else {
				results = append(results, result)
			}
		}
		for _, component := range buildPlan.Spec.Components.KustomizeBuildList {
			if result, err := component.Render(ctx, holos.PathComponent(instance.Dir)); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not render: %w", err))
			} else {
				results = append(results, result)
			}
		}
	}

	return results, nil
}

// findCueMod returns the root module location containing the cue.mod file or
// directory or an error if the builder arguments do not share a common root
// module.
func (b *Builder) findCueMod() (dir holos.PathCueMod, err error) {
	for _, origPath := range b.cfg.args {
		absPath, err := filepath.Abs(origPath)
		if err != nil {
			return "", err
		}
		path := holos.PathCueMod(absPath)
		for {
			if _, err := os.Stat(filepath.Join(string(path), "cue.mod")); err == nil {
				if dir != "" && dir != path {
					return "", fmt.Errorf("multiple modules not supported: %v is not %v", dir, path)
				}
				dir = path
				break
			} else if !os.IsNotExist(err) {
				return "", err
			}
			parentPath := holos.PathCueMod(filepath.Dir(string(path)))
			if parentPath == path {
				return "", fmt.Errorf("no cue.mod from root to leaf: %v", origPath)
			}
			path = parentPath
		}
	}
	return dir, nil
}
