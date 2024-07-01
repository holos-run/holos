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
	"strings"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/holos-run/holos/api/core/v1alpha2"
	"github.com/holos-run/holos/api/v1alpha1"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/client"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/render"
)

const (
	KubernetesObjects = v1alpha2.KubernetesObjectsKind
	// Helm is the value of the kind field of holos build output indicating helm
	// values and helm command information.
	Helm = v1alpha2.HelmChartKind
	// Skip is the value when the instance should be skipped
	Skip = "Skip"
	// KustomizeBuild is the value of the kind field of cue output indicating
	// holos should process the component using kustomize build to render output.
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

type buildPlanWrapper struct {
	buildPlan *v1alpha2.BuildPlan
}

func (b *buildPlanWrapper) validate() error {
	if b == nil {
		return fmt.Errorf("invalid BuildPlan: is nil")
	}
	bp := b.buildPlan
	if bp == nil {
		return fmt.Errorf("invalid BuildPlan: is nil")
	}
	errs := make([]string, 0, 2)
	if bp.Kind != v1alpha2.BuildPlanKind {
		errs = append(errs, fmt.Sprintf("kind invalid: want: %s have: %s", v1alpha1.BuildPlanKind, bp.Kind))
	}
	if bp.APIVersion != v1alpha2.APIVersion {
		errs = append(errs, fmt.Sprintf("apiVersion invalid: want: %s have: %s", v1alpha2.APIVersion, bp.APIVersion))
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid BuildPlan: " + strings.Join(errs, ", "))
	}
	return nil
}

func (b *buildPlanWrapper) resultCapacity() (count int) {
	if b == nil {
		return 0
	}
	bp := b.buildPlan
	count = len(bp.Spec.Components.HelmChartList) +
		len(bp.Spec.Components.KubernetesObjectsList) +
		len(bp.Spec.Components.KustomizeBuildList) +
		len(bp.Spec.Components.Resources)
	return count
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
func (b *Builder) Instances(ctx context.Context, cfg *client.Config) ([]*build.Instance, error) {
	log := logger.FromContext(ctx)

	mod, err := b.findCueMod()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	dir := string(mod)

	cueConfig := load.Config{Dir: dir}

	// Get the platform model from the PlatformConfig
	pc, err := client.LoadPlatformConfig(ctx, dir)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	data, err := json.Marshal(pc)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// Refer to https://github.com/cue-lang/cue/blob/v0.7.0/cmd/cue/cmd/common.go#L429
	cueConfig.Tags = append(cueConfig.Tags, "platform_config="+string(data))
	if b.Cluster() != "" {
		cueConfig.Tags = append(cueConfig.Tags, "cluster="+b.Cluster())
	}
	log.DebugContext(ctx, fmt.Sprintf("cue: tags %v", cueConfig.Tags))

	prefix := []string{"cue", "export", "--out", "yaml"}
	for _, tag := range cueConfig.Tags {
		prefix = append(prefix, "-t", fmt.Sprintf("'%s'", tag))
	}

	// Make args relative to the module directory
	args := make([]string, len(b.cfg.args))
	for idx, path := range b.cfg.args {
		target, err := filepath.Abs(path)
		if err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not find absolute path: %w", err))
		}
		relPath, err := filepath.Rel(dir, target)
		if err != nil {
			return nil, errors.Wrap(fmt.Errorf("invalid argument, must be relative to cue.mod: %w", err))
		}
		relPath = "./" + relPath
		args[idx] = relPath

		equiv := make([]string, len(prefix), 1+len(prefix))
		copy(equiv, prefix)
		equiv = append(equiv, relPath)
		log.Debug(strings.Join(equiv, " "), "comment", "cue equivalent command")
	}

	return load.Instances(args, &cueConfig), nil
}

func (b *Builder) Run(ctx context.Context, cfg *client.Config) (results []*render.Result, err error) {
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, "cue: building instances")
	instances, err := b.Instances(ctx, cfg)
	if err != nil {
		return nil, err
	}
	results = make([]*render.Result, 0, len(instances)*8)

	// Each CUE instance provides a BuildPlan
	for idx, instance := range instances {
		log.DebugContext(ctx, "cue: building instance", "idx", idx, "dir", instance.Dir)
		r, err := b.runInstance(ctx, instance)
		if err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not run: %w", err))
		}
		results = append(results, r...)
	}

	return results, nil
}

func (b Builder) runInstance(ctx context.Context, instance *build.Instance) (results []*render.Result, err error) {
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

	log.DebugContext(ctx, "cue: decoding holos build plan")
	jsonBytes, err := value.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not marshal cue instance %s: %w", instance.Dir, err))
	}
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	// Discriminate the type of build plan.
	tm := &v1alpha1.TypeMeta{}
	err = decoder.Decode(tm)
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("invalid BuildPlan: %s: %w", instance.Dir, err))
	}

	log.DebugContext(ctx, "cue: discriminated build kind: "+tm.Kind, "kind", tm.Kind, "apiVersion", tm.APIVersion)

	// New decoder for the full object
	decoder = json.NewDecoder(bytes.NewReader(jsonBytes))

	// TODO: When we release v1, explicitly allow unknown fields so we can add
	// fields without needing to bump the major version.  Disallow until we reach
	// v1 for clear error reporting.
	decoder.DisallowUnknownFields()

	switch tm.Kind {
	case "BuildPlan":
		var bp v1alpha2.BuildPlan
		if err = decoder.Decode(&bp); err != nil {
			err = errors.Wrap(fmt.Errorf("could not decode BuildPlan %s: %w", instance.Dir, err))
			return
		}
		results, err = b.buildPlan(ctx, &bp, path)
		if err != nil {
			return results, err
		}
	default:
		err = errors.Wrap(fmt.Errorf("unknown kind: %v", tm.Kind))
	}

	return results, err
}

func (b *Builder) buildPlan(ctx context.Context, buildPlan *v1alpha2.BuildPlan, path holos.InstancePath) (results []*render.Result, err error) {
	log := logger.FromContext(ctx)

	bpw := buildPlanWrapper{buildPlan: buildPlan}

	if err := bpw.validate(); err != nil {
		log.WarnContext(ctx, "could not validate", "skipped", true, "err", err)
		return nil, errors.Wrap(fmt.Errorf("could not validate %w", err))
	}

	if buildPlan.Spec.Disabled {
		log.DebugContext(ctx, "skipped: spec.disabled is true", "skipped", true)
		return
	}

	results = make([]*render.Result, 0, bpw.resultCapacity())
	log.DebugContext(ctx, "allocated results slice", "cap", bpw.resultCapacity())

	for _, component := range buildPlan.Spec.Components.Resources {
		ko := render.KubernetesObjects{Component: component}
		if result, err := ko.Render(ctx, path); err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not render: %w", err))
		} else {
			results = append(results, result)
		}
	}

	for _, component := range buildPlan.Spec.Components.KubernetesObjectsList {
		ko := render.KubernetesObjects{Component: component}
		if result, err := ko.Render(ctx, path); err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not render: %w", err))
		} else {
			results = append(results, result)
		}
	}
	for _, component := range buildPlan.Spec.Components.HelmChartList {
		hc := render.HelmChart{Component: component}
		if result, err := hc.Render(ctx, path); err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not render: %w", err))
		} else {
			results = append(results, result)
		}
	}
	for _, component := range buildPlan.Spec.Components.KustomizeBuildList {
		kb := render.KustomizeBuild{Component: component}
		if result, err := kb.Render(ctx, path); err != nil {
			return nil, errors.Wrap(fmt.Errorf("could not render: %w", err))
		} else {
			results = append(results, result)
		}
	}

	log.DebugContext(ctx, "returning results", "len", len(results))

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
