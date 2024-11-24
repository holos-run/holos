package v1alpha5

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"cuelang.org/go/cue"
	core "github.com/holos-run/holos/api/core/v1alpha5"
	"github.com/holos-run/holos/internal/artifact"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

// Platform represents a platform builder.
type Platform struct {
	Platform core.Platform
}

// Load loads from a cue value.
func (p *Platform) Load(v cue.Value) error {
	return errors.Wrap(v.Decode(&p.Platform))
}

func (p *Platform) Export(encoder holos.Encoder) error {
	if err := encoder.Encode(&p.Platform); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (p *Platform) Select(selectors ...holos.Selector) []holos.Component {
	components := make([]holos.Component, 0, len(p.Platform.Spec.Components))
	for _, component := range p.Platform.Spec.Components {
		if holos.IsSelected(component.Labels, selectors...) {
			components = append(components, &Component{component})
		}
	}
	return components
}

type Component struct {
	Component core.Component
}

func (c *Component) Describe() string {
	if val, ok := c.Component.Annotations["app.holos.run/description"]; ok {
		return val
	}
	return c.Component.Name
}

func (c *Component) Tags() ([]string, error) {
	size := 2 +
		len(c.Component.Parameters) +
		len(c.Component.Labels) +
		len(c.Component.Annotations)

	tags := make([]string, 0, size)
	for k, v := range c.Component.Parameters {
		tags = append(tags, k+"="+v)
	}
	// Inject holos component metadata tags.
	tags = append(tags, "holos_component_name="+c.Component.Name)
	tags = append(tags, "holos_component_path="+c.Component.Path)

	if len(c.Component.Labels) > 0 {
		labels, err := json.Marshal(c.Component.Labels)
		if err != nil {
			return nil, err
		}
		tags = append(tags, "holos_component_labels="+string(labels))
	}

	if len(c.Component.Annotations) > 0 {
		annotations, err := json.Marshal(c.Component.Annotations)
		if err != nil {
			return nil, err
		}
		tags = append(tags, "holos_component_annotations="+string(annotations))
	}

	return tags, nil
}

func (c *Component) WriteTo() string {
	return c.Component.WriteTo
}

func (c *Component) Labels() holos.Labels {
	return c.Component.Labels
}

func (c *Component) Path() string {
	return util.DotSlash(c.Component.Path)
}

var _ holos.BuildPlan = &BuildPlan{}

// BuildPlan represents a component builder.
type BuildPlan struct {
	core.BuildPlan
	Opts holos.BuildOpts
}

// Build builds a BuildPlan into Artifact files.
func (b *BuildPlan) Build(ctx context.Context) error {
	name := b.BuildPlan.Metadata.Name
	path := b.Opts.Path
	log := logger.FromContext(ctx).With("name", name, "path", path)
	msg := fmt.Sprintf("could not build %s", name)
	if b.BuildPlan.Spec.Disabled {
		log.WarnContext(ctx, fmt.Sprintf("%s: disabled", msg))
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	// One more for the producer
	g.SetLimit(b.Opts.Concurrency + 1)

	// Producer.
	g.Go(func() error {
		for _, a := range b.BuildPlan.Spec.Artifacts {
			msg := fmt.Sprintf("%s artifact %s", msg, a.Artifact)
			log := log.With("artifact", a.Artifact)
			if a.Skip {
				log.WarnContext(ctx, fmt.Sprintf("%s: skipped field is true", msg))
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// https://golang.org/doc/faq#closures_and_goroutines
				a := a
				// Worker.  Blocks if limit has been reached.
				g.Go(func() error {
					for _, gen := range a.Generators {
						switch gen.Kind {
						case "Resources":
							if err := b.resources(log, gen, b.Opts.Store); err != nil {
								return errors.Format("could not generate resources: %w", err)
							}
						case "Helm":
							if err := b.helm(ctx, log, gen, b.Opts.Store); err != nil {
								return errors.Format("could not generate helm: %w", err)
							}
						case "File":
							if err := b.file(log, gen, b.Opts.Store); err != nil {
								return errors.Format("could not generate file: %w", err)
							}
						default:
							return errors.Format("%s: unsupported kind %s", msg, gen.Kind)
						}
					}

					for _, t := range a.Transformers {
						switch t.Kind {
						case "Kustomize":
							if err := b.kustomize(ctx, log, t, b.Opts.Store); err != nil {
								return errors.Wrap(err)
							}
						case "Join":
							s := make([][]byte, 0, len(t.Inputs))
							for _, input := range t.Inputs {
								if data, ok := b.Opts.Store.Get(string(input)); ok {
									s = append(s, data)
								} else {
									return errors.Format("%s: missing %s", msg, input)
								}
							}
							data := bytes.Join(s, []byte(t.Join.Separator))
							if err := b.Opts.Store.Set(string(t.Output), data); err != nil {
								return errors.Format("%s: %w", msg, err)
							}
							log.Debug("set artifact: " + string(t.Output))
						default:
							return errors.Format("%s: unsupported kind %s", msg, t.Kind)
						}
					}

					for _, validator := range a.Validators {
						switch validator.Kind {
						case "Command":
							if err := b.validate(ctx, log, validator, b.Opts.Store); err != nil {
								return errors.Wrap(err)
							}
						default:
							return errors.Format("%s: unsupported kind %s", msg, validator.Kind)
						}
					}

					// Write the final artifact
					if err := b.Opts.Store.Save(b.Opts.WriteTo, string(a.Artifact)); err != nil {
						return errors.Format("%s: %w", msg, err)
					}
					log.DebugContext(ctx, "wrote "+filepath.Join(b.Opts.WriteTo, string(a.Artifact)))

					return nil
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	return g.Wait()
}

func (b *BuildPlan) Export(idx int, encoder holos.OrderedEncoder) error {
	if err := encoder.Encode(idx, &b.BuildPlan); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (b *BuildPlan) Load(v cue.Value) error {
	return errors.Wrap(v.Decode(&b.BuildPlan))
}

func (b *BuildPlan) file(
	log *slog.Logger,
	g core.Generator,
	store artifact.Store,
) error {
	data, err := os.ReadFile(filepath.Join(string(b.Opts.Path), string(g.File.Source)))
	if err != nil {
		return errors.Wrap(err)
	}
	if err := store.Set(string(g.Output), data); err != nil {
		return errors.Wrap(err)
	}
	log.Debug("set artifact: " + string(g.Output))
	return nil
}

func (b *BuildPlan) helm(
	ctx context.Context,
	log *slog.Logger,
	g core.Generator,
	store artifact.Store,
) error {
	chartName := g.Helm.Chart.Name
	log = log.With("chart", chartName)
	// Unnecessary? cargo cult copied from internal/cli/render/render.go
	if chartName == "" {
		return errors.New("missing chart name")
	}

	// Cache the chart by version to pull new versions. (#273)
	cacheDir := filepath.Join(string(b.Opts.Path), "vendor", g.Helm.Chart.Version)
	cachePath := filepath.Join(cacheDir, filepath.Base(chartName))

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		timeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		err := onceWithLock(log, timeout, cachePath, func() error {
			return b.cacheChart(ctx, log, cacheDir, g.Helm.Chart)
		})
		if err != nil {
			return errors.Format("could not cache chart: %w", err)
		}
	}

	// Write values file
	tempDir, err := os.MkdirTemp("", "holos.helm")
	if err != nil {
		return errors.Format("could not make temp dir: %w", err)
	}
	defer util.Remove(ctx, tempDir)

	data, err := yaml.Marshal(g.Helm.Values)
	if err != nil {
		return errors.Format("could not marshal values: %w", err)
	}

	valuesPath := filepath.Join(tempDir, "values.yaml")
	if err := os.WriteFile(valuesPath, data, 0666); err != nil {
		return errors.Wrap(fmt.Errorf("could not write values: %w", err))
	}
	log.DebugContext(ctx, "wrote"+valuesPath)

	// Run charts
	args := []string{"template"}
	if !g.Helm.EnableHooks {
		args = append(args, "--no-hooks")
	}
	for _, apiVersion := range g.Helm.APIVersions {
		args = append(args, "--api-versions", apiVersion)
	}
	if kubeVersion := g.Helm.KubeVersion; kubeVersion != "" {
		args = append(args, "--kube-version", kubeVersion)
	}
	args = append(args,
		"--include-crds",
		"--values", valuesPath,
		"--namespace", g.Helm.Namespace,
		"--kubeconfig", "/dev/null",
		"--version", g.Helm.Chart.Version,
		g.Helm.Chart.Release,
		cachePath,
	)
	helmOut, err := util.RunCmd(ctx, "helm", args...)
	if err != nil {
		stderr := helmOut.Stderr.String()
		lines := strings.Split(stderr, "\n")
		for _, line := range lines {
			log.DebugContext(ctx, line)
			if strings.HasPrefix(line, "Error:") {
				err = fmt.Errorf("%s: %w", line, err)
			}
		}
		return errors.Format("could not run helm template: %w", err)
	}

	// Set the artifact
	if err := store.Set(string(g.Output), helmOut.Stdout.Bytes()); err != nil {
		return errors.Format("could not store helm output: %w", err)
	}
	log.Debug("set artifact: " + string(g.Output))

	return nil
}

func (b *BuildPlan) resources(
	log *slog.Logger,
	g core.Generator,
	store artifact.Store,
) error {
	var size int
	for _, m := range g.Resources {
		size += len(m)
	}
	list := make([]core.Resource, 0, size)

	for _, m := range g.Resources {
		for _, r := range m {
			list = append(list, r)
		}
	}

	msg := fmt.Sprintf(
		"could not generate %s for %s path %s",
		g.Output,
		b.BuildPlan.Metadata.Name,
		b.Opts.Path,
	)

	buf, err := marshal(list)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	if err := store.Set(string(g.Output), buf.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	log.Debug("set artifact " + string(g.Output))
	return nil
}

func (b *BuildPlan) validate(
	ctx context.Context,
	log *slog.Logger,
	validator core.Validator,
	store artifact.Store,
) error {
	tempDir, err := os.MkdirTemp("", "holos.validate")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)
	msg := fmt.Sprintf(
		"could not validate %s path %s",
		b.BuildPlan.Metadata.Name,
		b.Opts.Path,
	)

	// Write the inputs
	for _, input := range validator.Inputs {
		path := string(input)
		if err := store.Save(tempDir, path); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		log.DebugContext(ctx, "wrote "+filepath.Join(tempDir, path))
	}

	if len(validator.Command.Args) < 1 {
		return errors.Format("%s: command args length must be at least 1", msg)
	}
	size := len(validator.Command.Args) + len(validator.Inputs)
	args := make([]string, 0, size)
	args = append(args, validator.Command.Args...)
	for _, input := range validator.Inputs {
		args = append(args, filepath.Join(tempDir, string(input)))
	}

	// Execute the validator
	if _, err = util.RunCmdA(ctx, b.Opts.Stderr, args[0], args[1:]...); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	return nil
}

func (b *BuildPlan) kustomize(
	ctx context.Context,
	log *slog.Logger,
	t core.Transformer,
	store artifact.Store,
) error {
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)
	msg := fmt.Sprintf(
		"could not transform %s for %s path %s",
		t.Output,
		b.BuildPlan.Metadata.Name,
		b.Opts.Path,
	)

	// Write the kustomization
	data, err := yaml.Marshal(t.Kustomize.Kustomization)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	path := filepath.Join(tempDir, "kustomization.yaml")
	if err := os.WriteFile(path, data, 0666); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	log.DebugContext(ctx, "wrote "+path)

	// Write the inputs
	for _, input := range t.Inputs {
		path := string(input)
		if err := store.Save(tempDir, path); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		log.DebugContext(ctx, "wrote "+filepath.Join(tempDir, path))
	}

	// Execute kustomize
	r, err := util.RunCmdW(ctx, b.Opts.Stderr, "kubectl", "kustomize", tempDir)
	if err != nil {
		kErr := r.Stderr.String()
		err = errors.Format("%s: could not run kustomize: %w", msg, err)
		log.ErrorContext(ctx, fmt.Sprintf("%s: stderr:\n%s", err.Error(), kErr), "err", err, "stderr", kErr)
		return err
	}

	// Store the artifact
	if err := store.Set(string(t.Output), r.Stdout.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	log.Debug("set artifact " + string(t.Output))

	return nil
}

func marshal(list []core.Resource) (buf bytes.Buffer, err error) {
	encoder := yaml.NewEncoder(&buf)
	defer encoder.Close()
	for _, item := range list {
		if err = encoder.Encode(item); err != nil {
			err = errors.Wrap(err)
			return
		}
	}
	return
}

// cacheChart stores a cached copy of Chart in the chart subdirectory of path.
//
// We assume the only method responsible for writing to chartDir is cacheChart
// itself.  cacheChart runs concurrently when rendering a platform.
//
// We rely on the atomicity of moving temporary directories into place on the
// same filesystem via os.Rename. If a syscall.EEXIST error occurs during
// renaming, it indicates that the cached chart already exists, which is
// expected when this function is called concurrently.
//
// TODO(jeff): Break the dependency on v1alpha5, make it work across versions as
// a utility function.
func (b *BuildPlan) cacheChart(
	ctx context.Context,
	log *slog.Logger,
	cacheDir string,
	chart core.Chart,
) error {
	// Add repositories
	repo := chart.Repository
	stderr := b.Opts.Stderr
	if repo.URL == "" {
		// repo update not needed for oci charts so this is debug instead of warn.
		log.DebugContext(ctx, "skipped helm repo add and update: repo url is empty")
	} else {
		if _, err := util.RunCmdW(ctx, stderr, "helm", "repo", "add", repo.Name, repo.URL); err != nil {
			return errors.Format("could not run helm repo add: %w", err)
		}
		if _, err := util.RunCmdW(ctx, stderr, "helm", "repo", "update", repo.Name); err != nil {
			return errors.Format("could not run helm repo update: %w", err)
		}
	}

	cacheTemp, err := os.MkdirTemp(cacheDir, chart.Name)
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, cacheTemp)

	cn := chart.Name
	if chart.Repository.Name != "" {
		cn = fmt.Sprintf("%s/%s", chart.Repository.Name, chart.Name)
	}
	helmOut, err := util.RunCmdW(ctx, stderr, "helm", "pull", "--destination", cacheTemp, "--untar=true", "--version", chart.Version, cn)
	if err != nil {
		stderr := helmOut.Stderr.String()
		lines := strings.Split(stderr, "\n")
		for _, line := range lines {
			log.DebugContext(ctx, line)
			if strings.HasPrefix(line, "Error:") {
				err = fmt.Errorf("%s: %w", line, err)
			}
		}
		return errors.Format("could not run helm pull: %w", err)
	}
	log.Debug("helm pull", "stdout", helmOut.Stdout, "stderr", helmOut.Stderr)

	items, err := os.ReadDir(cacheTemp)
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not read directory: %w", err))
	}
	if len(items) != 1 {
		return errors.Format("want: exactly one item, have: %+v", items)
	}
	item := items[0]

	src := filepath.Join(cacheTemp, item.Name())
	dst := filepath.Join(cacheDir, chart.Name)
	if err := os.Rename(src, dst); err != nil {
		var linkErr *os.LinkError
		if errors.As(err, &linkErr) && errors.Is(linkErr.Err, syscall.EEXIST) {
			log.DebugContext(ctx, "cache already exists", "chart", chart.Name, "chart_version", chart.Version, "path", dst)
		} else {
			return errors.Wrap(fmt.Errorf("could not rename: %w", err))
		}
	} else {
		log.DebugContext(ctx, fmt.Sprintf("renamed %s to %s", src, dst), "src", src, "dst", dst)
	}

	log.InfoContext(ctx,
		fmt.Sprintf("cached %s %s", chart.Name, chart.Version),
		"chart", chart.Name,
		"chart_version", chart.Version,
		"path", dst,
	)

	return nil
}

// onceWithLock obtains a filesystem lock with mkdir, then executes fn.  If the
// lock is already locked, onceWithLock waits for it to be released then returns
// without calling fn.
func onceWithLock(log *slog.Logger, ctx context.Context, path string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return errors.Wrap(err)
	}

	// Obtain a lock with a timeout.
	lockDir := path + ".lock"
	log = log.With("lock", lockDir)

	err := os.Mkdir(lockDir, 0777)
	if err == nil {
		defer os.RemoveAll(lockDir)
		log.DebugContext(ctx, fmt.Sprintf("acquired %s", lockDir))
		if err := fn(); err != nil {
			return errors.Wrap(err)
		}
		log.DebugContext(ctx, fmt.Sprintf("released %s", lockDir))
		return nil
	}

	// Wait until the lock is released then return.
	if os.IsExist(err) {
		log.DebugContext(ctx, fmt.Sprintf("blocked %s", lockDir))
		stillBlocked := time.After(5 * time.Second)
		deadLocked := time.After(10 * time.Second)
		for {
			select {
			case <-stillBlocked:
				log.WarnContext(ctx, fmt.Sprintf("waiting for %s to be released", lockDir))
			case <-deadLocked:
				log.WarnContext(ctx, fmt.Sprintf("still waiting for %s to be released (dead lock?)", lockDir))
			case <-time.After(100 * time.Millisecond):
				if _, err := os.Stat(lockDir); os.IsNotExist(err) {
					log.DebugContext(ctx, fmt.Sprintf("unblocked %s", lockDir))
					return nil
				}
			case <-ctx.Done():
				return errors.Wrap(ctx.Err())
			}
		}
	}

	// Unexpected error
	return errors.Wrap(err)
}
