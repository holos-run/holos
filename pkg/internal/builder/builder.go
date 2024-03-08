// Package builder is responsible for building fully rendered kubernetes api
// objects from various input directories. A directory may contain a platform
// spec or a component spec.
package builder

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/util"
	"github.com/holos-run/holos/pkg/wrapper"
)

const (
	// Kube is the value of the kind field of holos build output indicating
	// kubernetes api objects.
	Kube = "KubernetesObjects"
	// Helm is the value of the kind field of holos build output indicating helm
	// values and helm command information.
	Helm = "HelmChart"
	// Skip is the value when the instance should be skipped
	Skip = "Skip"
	// ChartDir is the chart cache directory name.
	ChartDir = "vendor"
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

// apiObjectMap is the shape of marshalled api objects returned from cue to the
// holos cli. A map is used to improve the clarity of error messages from cue.
type apiObjectMap map[string]map[string]string

// fileContentMap is a map of file names to file contents.
type fileContentMap map[string]string

// Result is the build result for display or writing.
type Result struct {
	Metadata  Metadata `json:"metadata,omitempty"`
	KsContent string   `json:"ksContent,omitempty"`
	// APIObjectMap holds the marshalled representation of api objects.
	APIObjectMap      apiObjectMap `json:"apiObjectMap,omitempty"`
	accumulatedOutput string
	Skip              bool
	// KustomizeFiles holds the files for a kustomize kustomization directory.
	KustomizeFiles fileContentMap `json:"kustomizeFiles"`
	// ResourcesFile is the file name used for api objects in kustomization.yaml
	ResourcesFile string `json:"resourcesFile,omitempty"`
}

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Chart struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Release    string     `json:"release"`
	Repository Repository `json:"repository"`
}

// A HelmChart represents a helm command to provide chart values in order to render kubernetes api objects.
type HelmChart struct {
	APIVersion    string   `json:"apiVersion"`
	Kind          string   `json:"kind"`
	Metadata      Metadata `json:"metadata"`
	KsContent     string   `json:"ksContent"`
	Namespace     string   `json:"namespace"`
	Chart         Chart    `json:"chart"`
	ValuesContent string   `json:"valuesContent"`
	// APIObjectMap holds the marshalled representation of api objects.
	APIObjectMap apiObjectMap `json:"APIObjectMap"`
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

// AccumulatedOutput returns the accumulated rendered output.
func (r *Result) AccumulatedOutput() string {
	return r.accumulatedOutput
}

// addAPIObjects adds the overlay api objects to accumulatedOutput.
func (r *Result) addOverlayObjects(log *slog.Logger) {
	b := []byte(r.AccumulatedOutput())
	kinds := make([]string, 0, len(r.APIObjectMap))
	// Sort the keys
	for kind := range r.APIObjectMap {
		kinds = append(kinds, kind)
	}
	slices.Sort(kinds)

	for _, kind := range kinds {
		v := r.APIObjectMap[kind]
		// Sort the keys
		names := make([]string, 0, len(v))
		for name := range v {
			names = append(names, name)
		}
		slices.Sort(names)

		for _, name := range names {
			yamlString := v[name]
			log.Debug(fmt.Sprintf("%s/%s", kind, name), "kind", kind, "name", name)
			b = util.EnsureNewline(b)
			header := fmt.Sprintf("---\n# Source: CUE apiObjects.%s.%s\n", kind, name)
			b = append(b, []byte(header+yamlString)...)
			b = util.EnsureNewline(b)
		}
	}
	r.accumulatedOutput = string(b)
}

// kustomize replaces the final output with the output of kustomize build if the
func (r *Result) kustomize(ctx context.Context) error {
	log := logger.FromContext(ctx)
	if r.ResourcesFile == "" {
		log.DebugContext(ctx, "skipping kustomize: no resourcesFile")
		return nil
	}
	if len(r.KustomizeFiles) < 1 {
		log.DebugContext(ctx, "skipping kustomize: no kustomizeFiles")
		return nil
	}
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return wrapper.Wrap(err)
	}
	defer remove(ctx, tempDir)

	// Write the main api object resources file for kustomize.
	target := filepath.Join(tempDir, r.ResourcesFile)
	b := []byte(r.AccumulatedOutput())
	b = util.EnsureNewline(b)
	if err := os.WriteFile(target, b, 0644); err != nil {
		return wrapper.Wrap(fmt.Errorf("could not write resources: %w", err))
	}
	log.DebugContext(ctx, "wrote: "+target, "op", "write", "path", target, "bytes", len(b))

	// Write the kustomization tree, kustomization.yaml must be in this map for kustomize to work.
	for file, content := range r.KustomizeFiles {
		target := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return wrapper.Wrap(err)
		}
		b := []byte(content)
		b = util.EnsureNewline(b)
		if err := os.WriteFile(target, b, 0644); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not write: %w", err))
		}
		log.DebugContext(ctx, "wrote: "+target, "op", "write", "path", target, "bytes", len(b))
	}

	// Run kustomize.
	kOut, err := util.RunCmd(ctx, "kubectl", "kustomize", tempDir)
	if err != nil {
		log.ErrorContext(ctx, kOut.Stderr.String())
		return wrapper.Wrap(err)
	}
	// Replace the accumulated output
	r.accumulatedOutput = kOut.Stdout.String()
	return nil
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
	log.DebugContext(ctx, "out: wrote "+path, "action", "write", "path", path, "status", "ok")
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
		log.Debug("cue: equivalent command: " + equiv)
	}

	// Refer to https://github.com/cue-lang/cue/blob/v0.7.0/cmd/cue/cmd/common.go#L429
	cfg.Tags = append(cfg.Tags, "cluster="+b.Cluster())
	log.DebugContext(ctx, fmt.Sprintf("cue: tags %v", cfg.Tags))

	return load.Instances(args, &cfg), nil
}

func (b *Builder) Run(ctx context.Context) (results []*Result, err error) {
	results = make([]*Result, 0, len(b.cfg.args))
	cueCtx := cuecontext.New()
	logger.FromContext(ctx).DebugContext(ctx, "cue: building instances")
	instances, err := b.Instances(ctx)
	if err != nil {
		return results, err
	}

	for _, instance := range instances {
		var info buildInfo
		var result Result
		log := logger.FromContext(ctx).With("dir", instance.Dir)
		results = append(results, &result)
		if err := instance.Err; err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not load: %w", err))
		}
		log.DebugContext(ctx, "cue: building instance")
		value := cueCtx.BuildInstance(instance)
		if err := value.Err(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not build: %w", err))
		}
		log.DebugContext(ctx, "cue: validating instance")
		if err := value.Validate(); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not validate: %w", err))
		}
		log.DebugContext(ctx, "cue: decoding holos component build info")
		if err := value.Decode(&info); err != nil {
			return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
		}

		log.DebugContext(ctx, "cue: processing holos component kind "+info.Kind)
		switch kind := info.Kind; kind {
		case Skip:
			result.Skip = true
		case Kube:
			// CUE directly provides the kubernetes api objects in result.Content
			if err := value.Decode(&result); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
			}
			result.addOverlayObjects(log)
		case Helm:
			var helmChart HelmChart
			// First decode into the result.  Helm will populate the api objects later.
			if err := value.Decode(&result); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
			}
			// Decode again into the helm chart struct to get the values content to provide to helm.
			if err := value.Decode(&helmChart); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not decode: %w", err))
			}
			// runHelm populates result.Content from helm template output.
			if err := runHelm(ctx, &helmChart, &result, holos.PathComponent(instance.Dir)); err != nil {
				return nil, err
			}
			result.addOverlayObjects(log)
			if err := result.kustomize(ctx); err != nil {
				return nil, wrapper.Wrap(fmt.Errorf("could not kustomize: %w", err))
			}
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

// runHelm provides the values produced by CUE to helm template and returns
// the rendered kubernetes api objects in the result.
func runHelm(ctx context.Context, hc *HelmChart, r *Result, path holos.PathComponent) error {
	log := logger.FromContext(ctx).With("chart", hc.Chart.Name)
	if hc.Chart.Name == "" {
		log.WarnContext(ctx, "skipping helm: no chart name specified, use a different component type")
		return nil
	}

	cachedChartPath := filepath.Join(string(path), ChartDir, filepath.Base(hc.Chart.Name))
	if isNotExist(cachedChartPath) {
		// Add repositories
		repo := hc.Chart.Repository
		if repo.URL != "" {
			out, err := util.RunCmd(ctx, "helm", "repo", "add", repo.Name, repo.URL)
			if err != nil {
				log.ErrorContext(ctx, "could not run helm", "stderr", out.Stderr.String(), "stdout", out.Stdout.String())
				return wrapper.Wrap(fmt.Errorf("could not run helm repo add: %w", err))
			}
			// Update repository
			out, err = util.RunCmd(ctx, "helm", "repo", "update", repo.Name)
			if err != nil {
				log.ErrorContext(ctx, "could not run helm", "stderr", out.Stderr.String(), "stdout", out.Stdout.String())
				return wrapper.Wrap(fmt.Errorf("could not run helm repo update: %w", err))
			}
		} else {
			log.DebugContext(ctx, "no chart repository url proceeding assuming oci chart")
		}

		// Cache the chart
		if err := cacheChart(ctx, path, ChartDir, hc.Chart); err != nil {
			return fmt.Errorf("could not cache chart: %w", err)
		}
	}

	// Write values file
	tempDir, err := os.MkdirTemp("", "holos")
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not make temp dir: %w", err))
	}
	defer remove(ctx, tempDir)

	valuesPath := filepath.Join(tempDir, "values.yaml")
	if err := os.WriteFile(valuesPath, []byte(hc.ValuesContent), 0644); err != nil {
		return wrapper.Wrap(fmt.Errorf("could not write values: %w", err))
	}
	log.DebugContext(ctx, "helm: wrote values", "path", valuesPath, "bytes", len(hc.ValuesContent))

	// Run charts
	chart := hc.Chart
	helmOut, err := util.RunCmd(ctx, "helm", "template", "--include-crds", "--values", valuesPath, "--namespace", hc.Namespace, "--kubeconfig", "/dev/null", "--version", chart.Version, chart.Release, cachedChartPath)
	if err != nil {
		stderr := helmOut.Stderr.String()
		lines := strings.Split(stderr, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Error:") {
				err = fmt.Errorf("%s: %w", line, err)
			}
		}
		return wrapper.Wrap(fmt.Errorf("could not run helm template: %w", err))
	}

	r.accumulatedOutput = helmOut.Stdout.String()

	return nil
}

// remove cleans up path, useful for temporary directories.
func remove(ctx context.Context, path string) {
	log := logger.FromContext(ctx)
	if err := os.RemoveAll(path); err != nil {
		log.WarnContext(ctx, "tmp: could not remove", "err", err, "path", path)
	} else {
		log.DebugContext(ctx, "tmp: removed", "path", path)
	}
}

func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// cacheChart stores a cached copy of Chart in the chart subdirectory of path.
func cacheChart(ctx context.Context, path holos.PathComponent, chartDir string, chart Chart) error {
	log := logger.FromContext(ctx)

	cacheTemp, err := os.MkdirTemp(string(path), chartDir)
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not make temp dir: %w", err))
	}
	defer remove(ctx, cacheTemp)

	chartName := chart.Name
	if chart.Repository.Name != "" {
		chartName = fmt.Sprintf("%s/%s", chart.Repository.Name, chart.Name)
	}
	helmOut, err := util.RunCmd(ctx, "helm", "pull", "--destination", cacheTemp, "--untar=true", "--version", chart.Version, chartName)
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not run helm pull: %w", err))
	}
	log.Debug("helm pull", "stdout", helmOut.Stdout, "stderr", helmOut.Stderr)

	cachePath := filepath.Join(string(path), chartDir)
	if err := os.Rename(cacheTemp, cachePath); err != nil {
		return wrapper.Wrap(fmt.Errorf("could not rename: %w", err))
	}
	log.InfoContext(ctx, "cached", "chart", chart.Name, "version", chart.Version, "path", cachePath)

	return nil
}
