package v1alpha5

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"cuelang.org/go/cue"
	core "github.com/holos-run/holos/api/core/v1alpha5"
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
var _ task = generatorTask{}
var _ task = transformersTask{}
var _ task = validatorTask{}

type task interface {
	id() string
	run(ctx context.Context) error
}

type taskParams struct {
	taskName      string
	buildPlanName string
	opts          holos.BuildOpts
}

func (t taskParams) id() string {
	return fmt.Sprintf("%s:%s/%s", t.opts.Path, t.buildPlanName, t.taskName)
}

type generatorTask struct {
	taskParams
	generator core.Generator
	wg        *sync.WaitGroup
}

func (t generatorTask) run(ctx context.Context) error {
	defer t.wg.Done()
	msg := fmt.Sprintf("could not build %s", t.id())
	switch t.generator.Kind {
	case "Resources":
		if err := t.resources(); err != nil {
			return errors.Format("%s: could not generate resources: %w", msg, err)
		}
	case "Helm":
		if err := t.helm(ctx); err != nil {
			return errors.Format("%s: could not generate helm: %w", msg, err)
		}
	case "File":
		if err := t.file(); err != nil {
			return errors.Format("%s: could not generate file: %w", msg, err)
		}
	default:
		return errors.Format("%s: unsupported kind %s", msg, t.generator.Kind)
	}
	return nil
}

func (t generatorTask) file() error {
	data, err := os.ReadFile(filepath.Join(string(t.opts.Path), string(t.generator.File.Source)))
	if err != nil {
		return errors.Wrap(err)
	}
	if err := t.opts.Store.Set(string(t.generator.Output), data); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (t generatorTask) helm(ctx context.Context) error {
	chartName := t.generator.Helm.Chart.Name
	// Unnecessary? cargo cult copied from internal/cli/render/render.go
	if chartName == "" {
		return errors.New("missing chart name")
	}
	// Cache the chart by version to pull new versions. (#273)
	cacheDir := filepath.Join(string(t.opts.Path), "vendor", t.generator.Helm.Chart.Version)
	cachePath := filepath.Join(cacheDir, filepath.Base(chartName))

	log := logger.FromContext(ctx)

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		err := func() error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			defer cancel()
			return onceWithLock(log, ctx, cachePath, func() error {
				return cacheChart(ctx, cacheDir, t.generator.Helm.Chart, t.opts.Stderr)
			})
		}()
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

	data, err := yaml.Marshal(t.generator.Helm.Values)
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
	if !t.generator.Helm.EnableHooks {
		args = append(args, "--no-hooks")
	}
	for _, apiVersion := range t.generator.Helm.APIVersions {
		args = append(args, "--api-versions", apiVersion)
	}
	if kubeVersion := t.generator.Helm.KubeVersion; kubeVersion != "" {
		args = append(args, "--kube-version", kubeVersion)
	}
	args = append(args,
		"--include-crds",
		"--values", valuesPath,
		"--namespace", t.generator.Helm.Namespace,
		"--kubeconfig", "/dev/null",
		"--version", t.generator.Helm.Chart.Version,
		t.generator.Helm.Chart.Release,
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
	if err := t.opts.Store.Set(string(t.generator.Output), helmOut.Stdout.Bytes()); err != nil {
		return errors.Format("could not store helm output: %w", err)
	}
	log.Debug("set artifact: " + string(t.generator.Output))

	return nil
}

func (t generatorTask) resources() error {
	var size int
	for _, m := range t.generator.Resources {
		size += len(m)
	}
	list := make([]core.Resource, 0, size)

	for _, m := range t.generator.Resources {
		for _, r := range m {
			list = append(list, r)
		}
	}

	msg := fmt.Sprintf(
		"could not generate %s for %s",
		t.generator.Output,
		t.id(),
	)

	buf, err := marshal(list)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	if err := t.opts.Store.Set(string(t.generator.Output), buf.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	return nil
}

type transformersTask struct {
	taskParams
	transformers []core.Transformer
	wg           *sync.WaitGroup
}

func (t transformersTask) run(ctx context.Context) error {
	defer t.wg.Done()
	for idx, transformer := range t.transformers {
		msg := fmt.Sprintf("could not build %s/%d", t.id(), idx)
		switch transformer.Kind {
		case "Kustomize":
			if err := kustomize(ctx, transformer, t.taskParams); err != nil {
				return errors.Wrap(err)
			}
		case "Join":
			s := make([][]byte, 0, len(transformer.Inputs))
			for _, input := range transformer.Inputs {
				if data, ok := t.opts.Store.Get(string(input)); ok {
					s = append(s, data)
				} else {
					return errors.Format("%s: missing %s", msg, input)
				}
			}
			data := bytes.Join(s, []byte(transformer.Join.Separator))
			if err := t.opts.Store.Set(string(transformer.Output), data); err != nil {
				return errors.Format("%s: %w", msg, err)
			}
		default:
			return errors.Format("%s: unsupported kind %s", msg, transformer.Kind)
		}
	}
	return nil
}

type validatorTask struct {
	taskParams
	validator core.Validator
	wg        *sync.WaitGroup
}

func (t validatorTask) run(ctx context.Context) error {
	defer t.wg.Done()
	msg := fmt.Sprintf("could not validate %s", t.id())
	switch kind := t.validator.Kind; kind {
	case "Command":
		if err := validate(ctx, t.validator, t.taskParams); err != nil {
			return errors.Wrap(err)
		}
	default:
		return errors.Format("%s: unsupported kind %s", msg, kind)
	}
	return nil
}

func worker(ctx context.Context, idx int, tasks chan task) error {
	log := logger.FromContext(ctx).With("worker", idx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task, ok := <-tasks:
			if !ok {
				log.DebugContext(ctx, fmt.Sprintf("worker %d returning: tasks chan closed", idx))
				return nil
			}
			log.DebugContext(ctx, fmt.Sprintf("worker %d task %s starting", idx, task.id()))
			if err := task.run(ctx); err != nil {
				return errors.Wrap(err)
			}
			log.DebugContext(ctx, fmt.Sprintf("worker %d task %s finished ok", idx, task.id()))
		}
	}
}

func buildArtifact(ctx context.Context, idx int, artifact core.Artifact, tasks chan task, buildPlanName string, opts holos.BuildOpts) error {
	var wg sync.WaitGroup
	msg := fmt.Sprintf("could not build %s artifact %s", buildPlanName, artifact.Artifact)
	// Process Generators concurrently
	for gid, gen := range artifact.Generators {
		task := generatorTask{
			taskParams: taskParams{
				taskName:      fmt.Sprintf("artifact/%d/generator/%d", idx, gid),
				buildPlanName: buildPlanName,
				opts:          opts,
			},
			generator: gen,
			wg:        &wg,
		}
		wg.Add(1)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tasks <- task:
		}
	}
	wg.Wait()

	// Process Transformers sequentially
	task := transformersTask{
		taskParams: taskParams{
			taskName:      fmt.Sprintf("artifact/%d/transformers", idx),
			buildPlanName: buildPlanName,
			opts:          opts,
		},
		transformers: artifact.Transformers,
		wg:           &wg,
	}
	wg.Add(1)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case tasks <- task:
	}
	wg.Wait()

	// Process Validators concurrently
	for vid, val := range artifact.Validators {
		task := validatorTask{
			taskParams: taskParams{
				taskName:      fmt.Sprintf("artifact/%d/validator/%d", idx, vid),
				buildPlanName: buildPlanName,
				opts:          opts,
			},
			validator: val,
			wg:        &wg,
		}
		wg.Add(1)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tasks <- task:
		}
	}
	wg.Wait()

	// Write the final artifact
	out := string(artifact.Artifact)
	if err := opts.Store.Save(opts.WriteTo, out); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, fmt.Sprintf("wrote %s", filepath.Join(opts.WriteTo, out)))

	return nil
}

// BuildPlan represents a component builder.
type BuildPlan struct {
	core.BuildPlan
	Opts holos.BuildOpts
}

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
	tasks := make(chan task)

	// Start the worker pool.
	for idx := 0; idx < max(1, b.Opts.Concurrency); idx++ {
		g.Go(func() error {
			return worker(ctx, idx, tasks)
		})
	}

	// Start one producer that fans out to one pipeline per artifact.
	g.Go(func() error {
		// Close the tasks chan when the producer returns.
		defer func() {
			log.DebugContext(ctx, "producer returning: closing tasks chan")
			close(tasks)
		}()
		// Separate error group for producers.
		p, ctx := errgroup.WithContext(ctx)
		for idx, a := range b.BuildPlan.Spec.Artifacts {
			p.Go(func() error {
				return buildArtifact(ctx, idx, a, tasks, b.Metadata.Name, b.Opts)
			})
		}
		// Wait on producers to finish.
		return errors.Wrap(p.Wait())
	})

	// Wait on workers to finish.
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

func cacheChart(ctx context.Context, cacheDir string, chart core.Chart, stderr io.Writer) error {
	log := logger.FromContext(ctx)
	// Add repositories
	repo := chart.Repository
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

	// Support chart.Name = "oci:/ghcr.io/akuity/kargo-charts/kargo"
	chartBaseName := path.Base(chart.Name)

	cacheTemp, err := os.MkdirTemp(cacheDir, chartBaseName)
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
	dst := filepath.Join(cacheDir, chartBaseName)
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

func kustomize(ctx context.Context, t core.Transformer, p taskParams) error {
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)
	msg := fmt.Sprintf(
		"could not transform %s for %s path %s",
		t.Output,
		p.buildPlanName,
		p.opts.Path,
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

	// Write the inputs
	for _, input := range t.Inputs {
		path := string(input)
		if err := p.opts.Store.Save(tempDir, path); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
	}

	// Execute kustomize
	r, err := util.RunCmdW(ctx, p.opts.Stderr, "kubectl", "kustomize", tempDir)
	if err != nil {
		return errors.Format("%s: could not run kustomize: %w", msg, err)
	}

	// Store the artifact
	if err := p.opts.Store.Set(string(t.Output), r.Stdout.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	return nil
}

func validate(ctx context.Context, validator core.Validator, p taskParams) error {
	store := p.opts.Store
	tempDir, err := os.MkdirTemp("", "holos.validate")
	if err != nil {
		return errors.Wrap(err)
	}
	// defer util.Remove(ctx, tempDir)
	msg := fmt.Sprintf("could not validate %s", p.id())

	// Write the inputs
	for _, input := range validator.Inputs {
		path := string(input)
		if err := store.Save(tempDir, path); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
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
	if _, err = util.RunCmdA(ctx, p.opts.Stderr, args[0], args[1:]...); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	return nil
}
