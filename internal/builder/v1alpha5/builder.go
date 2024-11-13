package v1alpha5

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/holos-run/holos"
	h "github.com/holos-run/holos"
	core "github.com/holos-run/holos/api/core/v1alpha5"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

func Unify(cueCtx *cue.Context, path string, tags []string) (bd holos.BuildData, err error) {
	if bd.ModuleRoot, bd.Dir, err = util.FindRootLeaf(path); err != nil {
		return bd, errors.Wrap(err)
	}

	cueConfig := load.Config{
		Dir:        bd.ModuleRoot,
		ModuleRoot: bd.ModuleRoot,
		Tags:       tags,
	}
	args := []string{bd.Dir}
	instances := load.Instances(args, &cueConfig)
	v := cueCtx.BuildInstance(instances[0])
	if err = v.Err(); err != nil {
		return bd, errors.Wrap(err)
	}
	bd.Value = v

	return
}

func LoadPlatform(path string, tags []string) (*Platform, error) {
	bd, err := Unify(cuecontext.New(), path, tags)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	decoder, err := bd.Decoder()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var platform Platform
	if err := decoder.Decode(&platform.Platform); err != nil {
		return nil, errors.Wrap(err)
	}
	return &platform, nil
}

// Platform represents a platform builder.
type Platform struct {
	Platform    core.Platform
	Concurrency int
	Stderr      io.Writer
}

// Build builds a Platform by concurrently building a BuildPlan for each
// platform component.  No artifact files are written directly, only indirectly
// by rendering each component.
func (p *Platform) Build(ctx context.Context, _ h.ArtifactMap) error {
	parentStart := time.Now()
	components := p.Platform.Spec.Components
	total := len(components)
	g, ctx := errgroup.WithContext(ctx)
	// Limit the number of concurrent goroutines due to CUE memory usage concerns
	// while rendering components.  One more for the producer.
	g.SetLimit(p.Concurrency + 1)
	// Spawn a producer because g.Go() blocks when the group limit is reached.
	g.Go(func() error {
		for idx := range components {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Capture idx to avoid issues with closure.  Fixed in Go 1.22.
				idx := idx
				component := &components[idx]
				// Worker go routine.  Blocks if limit has been reached.
				g.Go(func() error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						start := time.Now()
						log := logger.FromContext(ctx).With(
							"name", component.Name,
							"path", component.Path,
							"num", idx+1,
							"total", total,
						)
						log.DebugContext(ctx, "render component")

						tags := make([]string, 0, 2+len(component.Parameters))
						tags = append(tags, "holos_component_name="+component.Name)
						tags = append(tags, "holos_component_path="+component.Path)
						for key, value := range component.Parameters {
							tags = append(tags, fmt.Sprintf("%s=%s", key, value))
						}

						// Execute a sub-process to limit CUE memory usage.
						args := make([]string, 0, 10)
						args = append(args,
							"render",
							"component",
						)
						for _, tag := range tags {
							args = append(args, "--inject", tag)
						}
						if component.WriteTo != "" {
							args = append(args, "--write-to", component.WriteTo)
						}
						args = append(args, component.Path)
						result, err := util.RunCmd(ctx, "holos", args...)
						// I've lost an hour+ digging into why I couldn't see log output
						// from sub-processes.  Make sure to surface at least stderr from
						// sub-processes.
						_, _ = io.Copy(p.Stderr, result.Stderr)
						if err != nil {
							return errors.Wrap(fmt.Errorf("could not render component: %w", err))
						}

						duration := time.Since(start)
						msg := fmt.Sprintf(
							"rendered %s in %s",
							component.Name,
							duration,
						)
						log.InfoContext(ctx, msg, "duration", duration)
						return nil
					}
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Since(parentStart)
	msg := fmt.Sprintf("rendered platform in %s", duration)
	logger.FromContext(ctx).InfoContext(ctx, msg, "duration", duration, "version", p.Platform.APIVersion)
	return nil
}

// BuildPlan represents a component builder.
type BuildPlan struct {
	BuildPlan   core.BuildPlan
	Concurrency int
	Stderr      io.Writer
	// WriteTo --write-to=deploy flag
	WriteTo string
	// Path represents the path to the component
	Path h.InstancePath
}

// Build builds a BuildPlan into Artifact files.
func (b *BuildPlan) Build(ctx context.Context, am h.ArtifactMap) error {
	name := b.BuildPlan.Metadata.Name
	path := b.BuildPlan.Source.Component.Path
	log := logger.FromContext(ctx).With("name", name, "path", path)
	msg := fmt.Sprintf("could not build %s", name)
	if b.BuildPlan.Spec.Disabled {
		log.WarnContext(ctx, fmt.Sprintf("%s: disabled", msg))
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	// One more for the producer
	g.SetLimit(b.Concurrency + 1)

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
							if err := b.resources(log, gen, am); err != nil {
								return errors.Format("could not generate resources: %w", err)
							}
						case "Helm":
							if err := b.helm(ctx, log, gen, am); err != nil {
								return errors.Format("could not generate helm: %w", err)
							}
						case "File":
							if err := b.file(log, gen, am); err != nil {
								return errors.Format("could not generate file: %w", err)
							}
						default:
							return errors.Format("%s: unsupported kind %s", msg, gen.Kind)
						}
					}

					for _, t := range a.Transformers {
						switch t.Kind {
						case "Kustomize":
							if err := b.kustomize(ctx, log, t, am); err != nil {
								return errors.Wrap(err)
							}
						case "Join":
							s := make([][]byte, 0, len(t.Inputs))
							for _, input := range t.Inputs {
								if data, ok := am.Get(string(input)); ok {
									s = append(s, data)
								} else {
									return errors.Format("%s: missing %s", msg, input)
								}
							}
							data := bytes.Join(s, []byte(t.Join.Separator))
							if err := am.Set(string(t.Output), data); err != nil {
								return errors.Format("%s: %w", msg, err)
							}
							log.Debug("set artifact: " + string(t.Output))
						default:
							return errors.Format("%s: unsupported kind %s", msg, t.Kind)
						}
					}

					// Write the final artifact
					if err := am.Save(b.WriteTo, string(a.Artifact)); err != nil {
						return errors.Format("%s: %w", msg, err)
					}
					log.DebugContext(ctx, "wrote "+filepath.Join(b.WriteTo, string(a.Artifact)))

					return nil
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	return g.Wait()
}

func (b *BuildPlan) file(
	log *slog.Logger,
	g core.Generator,
	am h.ArtifactMap,
) error {
	data, err := os.ReadFile(filepath.Join(string(b.Path), string(g.File.Source)))
	if err != nil {
		return errors.Wrap(err)
	}
	if err := am.Set(string(g.Output), data); err != nil {
		return errors.Wrap(err)
	}
	log.Debug("set artifact: " + string(g.Output))
	return nil
}

func (b *BuildPlan) helm(
	ctx context.Context,
	log *slog.Logger,
	g core.Generator,
	am h.ArtifactMap,
) error {
	chartName := g.Helm.Chart.Name
	log = log.With("chart", chartName)
	// Unnecessary? cargo cult copied from internal/cli/render/render.go
	if chartName == "" {
		return errors.New("missing chart name")
	}

	// Cache the chart by version to pull new versions. (#273)
	cacheDir := filepath.Join(string(b.Path), "vendor", g.Helm.Chart.Version)
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
			if strings.HasPrefix(line, "Error:") {
				err = fmt.Errorf("%s: %w", line, err)
			}
		}
		return errors.Format("could not run helm template: %w", err)
	}

	// Set the artifact
	if err := am.Set(string(g.Output), helmOut.Stdout.Bytes()); err != nil {
		return errors.Format("could not store helm output: %w", err)
	}
	log.Debug("set artifact: " + string(g.Output))

	return nil
}

func (b *BuildPlan) resources(
	log *slog.Logger,
	g core.Generator,
	am h.ArtifactMap,
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
		b.BuildPlan.Source.Component.Path,
	)

	buf, err := marshal(list)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	if err := am.Set(string(g.Output), buf.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	log.Debug("set artifact " + string(g.Output))
	return nil
}

func (b *BuildPlan) kustomize(
	ctx context.Context,
	log *slog.Logger,
	t core.Transformer,
	am h.ArtifactMap,
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
		b.BuildPlan.Source.Component.Path,
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
		if err := am.Save(tempDir, path); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		log.DebugContext(ctx, "wrote "+filepath.Join(tempDir, path))
	}

	// Execute kustomize
	r, err := util.RunCmd(ctx, "kubectl", "kustomize", tempDir)
	if err != nil {
		kErr := r.Stderr.String()
		err = errors.Format("%s: could not run kustomize: %w", msg, err)
		log.ErrorContext(ctx, fmt.Sprintf("%s: stderr:\n%s", err.Error(), kErr), "err", err, "stderr", kErr)
		return err
	}

	// Store the artifact
	if err := am.Set(string(t.Output), r.Stdout.Bytes()); err != nil {
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
	if repo.URL == "" {
		// repo update not needed for oci charts so this is debug instead of warn.
		log.DebugContext(ctx, "skipped helm repo add and update: repo url is empty")
	} else {
		if r, err := util.RunCmd(ctx, "helm", "repo", "add", repo.Name, repo.URL); err != nil {
			_, _ = io.Copy(b.Stderr, r.Stderr)
			return errors.Format("could not run helm repo add: %w", err)
		}
		if r, err := util.RunCmd(ctx, "helm", "repo", "update", repo.Name); err != nil {
			_, _ = io.Copy(b.Stderr, r.Stderr)
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
	helmOut, err := util.RunCmd(ctx, "helm", "pull", "--destination", cacheTemp, "--untar=true", "--version", chart.Version, cn)
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not run helm pull: %w", err))
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
