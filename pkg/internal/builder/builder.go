// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"context"
	"cuelang.org/go/cue/build"
	"fmt"
	"github.com/holos-run/holos"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/wrapper"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

const (
	// Kube is the value of the kind field of holos build output indicating
	// kubernetes api objects.
	Kube = "KubernetesObjects"
	// Helm is the value of the kind field of holos build output indicating helm
	// values and helm command information.
	Helm = "ChartValues"
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

type buildInfo struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

// Metadata represents the standard metadata fields of the cue output
type Metadata struct {
	Name string `json:"name,omitempty"`
}

// Result is the build result for display or writing.
type Result struct {
	Metadata  Metadata `json:"metadata,omitempty"`
	Content   string   `json:"content,omitempty"`
	KsContent string   `json:"ksContent,omitempty"`
}

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Chart struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository string `json:"repository"`
}

// A ChartValues represents a helm command to provide chart values in order to render kubernetes api objects.
type ChartValues struct {
	APIVersion   string       `json:"apiVersion"`
	Kind         string       `json:"kind"`
	Metadata     Metadata     `json:"metadata"`
	KsContent    string       `json:"ksContent"`
	Repositories []Repository `json:"repositories"`
	Charts       []Chart      `json:"charts"`
}

// Name returns the metadata name of the result. Equivalent to the
// OrderedComponent name specified in platform.yaml in the holos prototype.
func (r *Result) Name() string {
	return r.Metadata.Name
}

func (r *Result) Filename(writeTo string, cluster string) string {
	return filepath.Join(writeTo, "clusters", cluster, "components", r.Name(), r.Name()+".gen.yaml")
}

func (r *Result) KustomizationFilename(writeTo string, cluster string) string {
	return filepath.Join(writeTo, "clusters", cluster, "holos", "components", r.Name()+"-kustomization.gen.yaml")
}

// Save writes the content to the filesystem for git ops.
func (r *Result) Save(ctx context.Context, path string, content string) error {
	log := logger.FromContext(ctx)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.FileMode(0775)); err != nil {
		log.WarnContext(ctx, "could not mkdir", "path", dir, "err", err)
		return wrapper.Wrap(err)
	}
	// Write the kube api objects
	if err := os.WriteFile(path, []byte(content), os.FileMode(0644)); err != nil {
		log.WarnContext(ctx, "could not write", "path", path, "err", err)
		return wrapper.Wrap(err)
	}
	log.DebugContext(ctx, "wrote "+path, "action", "write", "path", path, "status", "ok")
	return nil
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
		log.Debug(equiv)
	}

	// Refer to https://github.com/cue-lang/cue/blob/v0.7.0/cmd/cue/cmd/common.go#L429
	cfg.Tags = append(cfg.Tags, "cluster="+b.Cluster())
	log.DebugContext(ctx, fmt.Sprintf("configured cue tags: %v", cfg.Tags))

	return load.Instances(args, &cfg), nil
}

func (b *Builder) Run(ctx context.Context) (results []*Result, err error) {
	results = make([]*Result, 0, len(b.cfg.args))
	cueCtx := cuecontext.New()
	instances, err := b.Instances(ctx)
	if err != nil {
		return results, err
	}

	for _, instance := range instances {
		var info buildInfo
		var result Result
		results = append(results, &result)
		if err := instance.Err; err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not load: %w", err))
		}
		value := cueCtx.BuildInstance(instance)
		if err := value.Err(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not build: %w", err))
		}
		if err := value.Validate(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not validate: %w", err))
		}

		if err := value.Decode(&info); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
		}

		switch kind := info.Kind; kind {
		case Kube:
			// TODO: Decode into a intermediate struct
			if err := value.Decode(&result); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
			}
		case Helm:
			var chartValues ChartValues
			if err := value.Decode(&chartValues); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
			}
			fmt.Printf("%#v\n\n", chartValues)
			return nil, wrapper.Wrap(fmt.Errorf("helm not implemented"))
		default:
			return nil, wrapper.Wrap(fmt.Errorf("build kind not implemented: %v", kind))
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
