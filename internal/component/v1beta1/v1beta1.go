// Package v1beta1 executes one component's [core.TaskSet].  The executor
// derives DAG edges from task inputs/output declarations plus explicit
// dependsOn edges per doc/design/v1beta1/schema.md D1, then executes tasks in
// topological order with bounded concurrency.  Scope is intra-component only;
// the platform-wide DAG belongs to Phase 2 (doc/design/v1beta1/rendering.md).
package v1beta1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"cuelang.org/go/cue"
	core "github.com/holos-run/holos/api/core/v1beta1"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/helm"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/cli"
)

var _ holos.BuildPlan = &TaskSet{}

// taskNamePattern validates task names per schema.md D3: an RFC 1123 label.
var taskNamePattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// TaskSet represents a component builder executing one [core.TaskSet].
type TaskSet struct {
	core.TaskSet
	Opts holos.BuildOpts

	// runHook is a test seam wrapping task execution.  When nil, run is called
	// directly.  Tests may instrument or replace run to observe scheduling
	// behavior without executing task bodies.
	runHook func(ctx context.Context, name string, run func(context.Context) error) error

	// saveMu guards saved so concurrent tasks materialize each store path into
	// the shared build temp directory at most once.
	saveMu sync.Mutex
	saved  map[string]error
}

// sharedSave materializes a store path into the shared build temp directory at
// most once.  Store paths are write-once and immutable, so the first
// materialization is authoritative; concurrent tasks consuming the same path
// must not race by truncating and rewriting a file another task is reading.
func (b *TaskSet) sharedSave(dir, path string) error {
	b.saveMu.Lock()
	defer b.saveMu.Unlock()
	key := dir + "\x00" + path
	if err, ok := b.saved[key]; ok {
		return err
	}
	err := b.Opts.Store.Save(dir, path)
	b.saved[key] = err
	return err
}

// Load loads the TaskSet from a cue value.
func (b *TaskSet) Load(v cue.Value) error {
	// First validate the value to get better error messages
	if err := v.Validate(cue.Concrete(true)); err != nil {
		return err
	}

	if err := v.Decode(&b.TaskSet); err != nil {
		// If it's a CUE error, return it unwrapped to preserve CUE's error formatting
		if v.Err() != nil {
			return v.Err()
		}
		return errors.Wrap(err)
	}
	return nil
}

// Export encodes the TaskSet at index idx.
func (b *TaskSet) Export(idx int, encoder holos.OrderedEncoder) error {
	if err := encoder.Encode(idx, &b.TaskSet); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Build derives the task graph per schema.md D1, then executes tasks in
// topological order with concurrency bounded by Opts.Concurrency.
func (b *TaskSet) Build(ctx context.Context) error {
	name := b.Metadata.Name
	log := logger.FromContext(ctx).With(
		"name", name,
		"path", b.Opts.Leaf(),
	)

	msg := fmt.Sprintf("could not build %s", name)
	if b.Spec.Disabled {
		log.WarnContext(ctx, fmt.Sprintf("%s: disabled", msg))
		return nil
	}

	g, err := b.graph()
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	b.saveMu.Lock()
	b.saved = make(map[string]error)
	b.saveMu.Unlock()

	// Load inputs sourced from the component directory into the artifact
	// store so tasks consume them uniformly (schema.md D1: an input matching
	// no output must exist in the component directory).
	for _, path := range g.files {
		if err := b.Opts.Store.Load(b.Opts.AbsLeaf(), path); err != nil {
			return errors.Format("%s: could not load %s from component directory: %w", msg, path, err)
		}
	}

	if err := b.execute(ctx, g); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	return nil
}

// graph represents the task DAG derived per schema.md D1.
type graph struct {
	// names holds every task name sorted for deterministic iteration.
	names []string
	// succ maps a task to the set of tasks depending on it.
	succ map[string]map[string]struct{}
	// pred maps a task to the set of tasks it depends on.
	pred map[string]map[string]struct{}
	// files holds input paths read from the component directory, sorted.
	files []string
}

// addEdge adds edge from -> to, ignoring duplicates (a duplicate edge is
// harmless and legal per schema.md D1).
func (g *graph) addEdge(from, to string) {
	if _, ok := g.succ[from][to]; ok {
		return
	}
	g.succ[from][to] = struct{}{}
	g.pred[to][from] = struct{}{}
}

// graph validates the task set and derives DAG edges per schema.md D1:
// inputs/output matching plus explicit dependsOn edges, unioned.  Cycles and
// unknown references are errors naming the tasks involved.
func (b *TaskSet) graph() (*graph, error) {
	tasks := b.Spec.Tasks

	g := &graph{
		names: make([]string, 0, len(tasks)),
		succ:  make(map[string]map[string]struct{}, len(tasks)),
		pred:  make(map[string]map[string]struct{}, len(tasks)),
	}
	for name := range tasks {
		g.names = append(g.names, name)
		g.succ[name] = make(map[string]struct{})
		g.pred[name] = make(map[string]struct{})
	}
	sort.Strings(g.names)

	// Validate tasks and index producers.  Outputs are write-once: a second
	// task declaring an already-declared output is an error naming both tasks.
	// Final artifact paths are write-once for the same reason: two sinks
	// declaring the same path is an error naming both tasks.
	producers := make(map[string]string, len(tasks))
	artifactPaths := make(map[string]string)
	for _, name := range g.names {
		task := tasks[name]
		if err := validateTask(name, task); err != nil {
			return nil, err
		}
		if output := string(task.Output); output != "" {
			if prev, ok := producers[output]; ok {
				return nil, errors.Format("duplicate output %s: declared by tasks %s and %s", output, prev, name)
			}
			producers[output] = name
		}
		if task.Kind == "Artifact" {
			path := string(task.Artifact.Path)
			if path == "" {
				path = string(task.Inputs[0])
			}
			if prev, ok := artifactPaths[path]; ok {
				return nil, errors.Format("duplicate artifact path %s: declared by tasks %s and %s", path, prev, name)
			}
			artifactPaths[path] = name
		}
	}

	// Sorted outputs for deterministic prefix matching.
	outputs := make([]string, 0, len(producers))
	for output := range producers {
		outputs = append(outputs, output)
	}
	sort.Strings(outputs)

	files := make(map[string]struct{})
	for _, name := range g.names {
		task := tasks[name]
		for _, input := range task.Inputs {
			path := string(input)
			matches := matchProducers(path, producers, outputs)
			if len(matches) == 0 {
				// The input must exist in the component directory.
				if _, err := os.Stat(filepath.Join(b.Opts.AbsLeaf(), path)); err != nil {
					return nil, errors.Format("task %s: input %s matches no task output and does not exist in the component directory", name, path)
				}
				files[path] = struct{}{}
				continue
			}
			for _, producer := range matches {
				if producer == name {
					return nil, errors.Format("task %s: input %s matches its own output", name, path)
				}
				g.addEdge(producer, name)
			}
		}

		// dependsOn adds explicit edges by task name for ordering constraints
		// with no data flow.
		targets := make([]string, 0, len(task.DependsOn))
		for target := range task.DependsOn {
			targets = append(targets, target)
		}
		sort.Strings(targets)
		for _, target := range targets {
			if _, ok := tasks[target]; !ok {
				if strings.Contains(target, ":") {
					return nil, errors.Format("task %s: dependsOn target %s: canonical task ids are not supported by intra-component execution", name, target)
				}
				return nil, errors.Format("task %s: dependsOn target %s: no such task", name, target)
			}
			if target == name {
				return nil, errors.Format("task %s: dependsOn target %s: task depends on itself", name, target)
			}
			g.addEdge(target, name)
		}
	}

	g.files = make([]string, 0, len(files))
	for path := range files {
		g.files = append(g.files, path)
	}
	sort.Strings(g.files)

	if err := checkCycles(g); err != nil {
		return nil, err
	}

	return g, nil
}

// matchProducers finds the tasks producing input path per schema.md D1.
// Three rules are tried in order; the first rule yielding matches wins.
func matchProducers(path string, producers map[string]string, outputs []string) []string {
	// Rule 1: exact match.
	if producer, ok := producers[path]; ok {
		return []string{producer}
	}
	// Rule 2: directory input.  A directory input may legitimately match
	// several producers and gains an edge from each.
	var matches []string
	prefix := path + "/"
	for _, output := range outputs {
		if strings.HasPrefix(output, prefix) {
			matches = append(matches, producers[output])
		}
	}
	if len(matches) > 0 {
		return matches
	}
	// Rule 3: directory output.  The input falls under a produced directory.
	for _, output := range outputs {
		if strings.HasPrefix(path, output+"/") {
			matches = append(matches, producers[output])
		}
	}
	return matches
}

// validateTask revalidates the per-kind constraints of schema.md at execution
// time: task name pattern (D3), the inputs/output requiredness and cardinality
// table (Task kinds), and path containment.  Every declared path must stay
// local to its root directory: store paths resolve under the build temp dir,
// file sources under the component directory, and artifact paths under the
// write-to directory.
func validateTask(name string, task core.Task) error {
	if !taskNamePattern.MatchString(name) {
		return errors.Format("invalid task name %q: must match %s", name, taskNamePattern)
	}
	if output := string(task.Output); output != "" && !filepath.IsLocal(output) {
		return errors.Format("task %s: output %s: path must be relative and must not traverse outside the build directory", name, output)
	}
	for _, input := range task.Inputs {
		if !filepath.IsLocal(string(input)) {
			return errors.Format("task %s: input %s: path must be relative and must not traverse outside the build directory", name, input)
		}
	}
	switch task.Kind {
	case "Resources", "Helm", "File":
		if len(task.Inputs) != 0 {
			return errors.Format("task %s: kind %s must not declare inputs", name, task.Kind)
		}
		if task.Output == "" {
			return errors.Format("task %s: kind %s requires an output", name, task.Kind)
		}
		if task.Kind == "File" && !filepath.IsLocal(string(task.File.Source)) {
			return errors.Format("task %s: file source %s: path must be relative and must not traverse outside the component directory", name, task.File.Source)
		}
	case "Kustomize", "Join":
		if len(task.Inputs) < 1 {
			return errors.Format("task %s: kind %s requires at least one input", name, task.Kind)
		}
		if task.Output == "" {
			return errors.Format("task %s: kind %s requires an output", name, task.Kind)
		}
		for path := range task.Kustomize.Files {
			if !filepath.IsLocal(string(path)) {
				return errors.Format("task %s: kustomize file %s: path must be relative and must not traverse outside the kustomize directory", name, path)
			}
		}
	case "Command":
		if len(task.Command.Args) < 1 {
			return errors.Format("task %s: command args length must be at least 1", name)
		}
		if task.Command.IsStdoutOutput && task.Output == "" {
			return errors.Format("task %s: kind Command requires an output when isStdoutOutput is true", name)
		}
		if stdin := task.Command.Stdin; stdin != "" {
			if !slices.Contains(task.Inputs, stdin) {
				return errors.Format("task %s: command stdin %s must be one of the task inputs", name, stdin)
			}
		}
	case "Artifact":
		if len(task.Inputs) != 1 {
			return errors.Format("task %s: kind Artifact requires exactly one input", name)
		}
		if task.Output != "" {
			return errors.Format("task %s: kind Artifact must not declare an output", name)
		}
		if path := string(task.Artifact.Path); path != "" && !filepath.IsLocal(path) {
			return errors.Format("task %s: artifact path %s: path must be relative and must not traverse outside the write-to directory", name, path)
		}
	default:
		return errors.Format("task %s: unsupported kind %s", name, task.Kind)
	}
	return nil
}

// checkCycles runs Kahn's algorithm over the graph.  Any cycle is an error
// reported with the full cycle path per schema.md D1.
func checkCycles(g *graph) error {
	indegree := make(map[string]int, len(g.names))
	queue := make([]string, 0, len(g.names))
	for _, name := range g.names {
		indegree[name] = len(g.pred[name])
		if indegree[name] == 0 {
			queue = append(queue, name)
		}
	}
	processed := 0
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		processed++
		for succ := range g.succ[name] {
			indegree[succ]--
			if indegree[succ] == 0 {
				queue = append(queue, succ)
			}
		}
	}
	if processed == len(g.names) {
		return nil
	}

	// Every remaining task participates in or depends on a cycle.  Walk
	// predecessors from the smallest remaining task until one repeats.
	remaining := make(map[string]bool, len(g.names)-processed)
	var start string
	for _, name := range g.names {
		if indegree[name] > 0 {
			remaining[name] = true
			if start == "" {
				start = name
			}
		}
	}
	seen := make(map[string]int)
	path := make([]string, 0, len(remaining))
	cur := start
	for {
		if idx, ok := seen[cur]; ok {
			cycle := path[idx:]
			// Predecessor walk yields the cycle in reverse execution order.
			slices.Reverse(cycle)
			return errors.Format("cycle detected: %s -> %s", strings.Join(cycle, " -> "), cycle[0])
		}
		seen[cur] = len(path)
		path = append(path, cur)
		// Deterministically choose the smallest remaining predecessor.
		next := ""
		for pred := range g.pred[cur] {
			if remaining[pred] && (next == "" || pred < next) {
				next = pred
			}
		}
		cur = next
	}
}

// execute runs tasks in topological order.  Ready tasks run concurrently on
// an errgroup bounded by Opts.Concurrency.  The ready queue is kept sorted by
// task name so dispatch order is a deterministic function of completion order.
func (b *TaskSet) execute(ctx context.Context, g *graph) error {
	eg, egctx := errgroup.WithContext(ctx)
	eg.SetLimit(max(1, b.Opts.Concurrency))

	indegree := make(map[string]int, len(g.names))
	ready := make([]string, 0, len(g.names))
	for _, name := range g.names {
		indegree[name] = len(g.pred[name])
		if indegree[name] == 0 {
			ready = append(ready, name)
		}
	}

	// completions is buffered for every task so workers never block sending.
	completions := make(chan string, len(g.names))
	pending := len(g.names)

	for pending > 0 {
		// Dispatch every ready task in sorted order.  Go blocks when the
		// concurrency limit is reached until a worker returns.
		for len(ready) > 0 && egctx.Err() == nil {
			name := ready[0]
			ready = ready[1:]
			eg.Go(func() error {
				if err := egctx.Err(); err != nil {
					return err
				}
				if err := b.runTask(egctx, name); err != nil {
					return err
				}
				completions <- name
				return nil
			})
		}

		select {
		case name := <-completions:
			pending--
			for succ := range g.succ[name] {
				indegree[succ]--
				if indegree[succ] == 0 {
					idx, _ := slices.BinarySearch(ready, succ)
					ready = slices.Insert(ready, idx, succ)
				}
			}
		case <-egctx.Done():
			return eg.Wait()
		}
	}
	return eg.Wait()
}

// runTask executes one task by kind.
func (b *TaskSet) runTask(ctx context.Context, name string) error {
	t := &taskRunner{
		name:        name,
		taskSetName: b.Metadata.Name,
		task:        b.Spec.Tasks[name],
		opts:        b.Opts,
		sharedSave:  b.sharedSave,
	}
	if b.runHook != nil {
		return b.runHook(ctx, name, t.run)
	}
	return t.run(ctx)
}

// taskRunner executes one [core.Task].
type taskRunner struct {
	name        string
	taskSetName string
	task        core.Task
	opts        holos.BuildOpts
	// sharedSave materializes a store path into the shared build temp
	// directory at most once across concurrent tasks.
	sharedSave func(dir, path string) error
}

// id uniquely identifies the task for log and error messages.
func (t *taskRunner) id() string {
	return fmt.Sprintf("%s:%s/%s", t.opts.Leaf(), t.taskSetName, t.name)
}

func (t *taskRunner) tempDir() (string, error) {
	if tempDir := t.opts.TempDir(); tempDir == "" {
		return "", errors.Format("missing build context temp directory")
	} else {
		return tempDir, nil
	}
}

func (t *taskRunner) run(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.DebugContext(ctx, fmt.Sprintf("task %s starting", t.id()))
	msg := fmt.Sprintf("could not build %s", t.id())
	switch t.task.Kind {
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
	case "Kustomize":
		if err := t.kustomize(ctx); err != nil {
			return errors.Format("%s: could not kustomize: %w", msg, err)
		}
	case "Join":
		if err := t.join(); err != nil {
			return errors.Format("%s: could not join: %w", msg, err)
		}
	case "Command":
		if err := t.command(ctx); err != nil {
			return errors.Format("%s: could not run command: %w", msg, err)
		}
	case "Artifact":
		if err := t.artifact(ctx); err != nil {
			return errors.Format("%s: could not write artifact: %w", msg, err)
		}
	default:
		return errors.Format("%s: unsupported kind %s", msg, t.task.Kind)
	}
	log.DebugContext(ctx, fmt.Sprintf("task %s finished ok", t.id()))
	return nil
}

// resources marshals kubernetes resources defined in CUE into the output.
func (t *taskRunner) resources() error {
	var size int
	for _, m := range t.task.Resources {
		size += len(m)
	}
	list := make([]core.Resource, 0, size)

	for _, m := range t.task.Resources {
		for _, r := range m {
			list = append(list, r)
		}
	}

	msg := fmt.Sprintf("could not generate %s for %s", t.task.Output, t.id())

	buf, err := marshal(list)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	if err := t.opts.Store.Set(string(t.task.Output), buf.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	return nil
}

// file reads a file from the component directory into the output.
func (t *taskRunner) file() error {
	data, err := os.ReadFile(filepath.Join(t.opts.AbsLeaf(), string(t.task.File.Source)))
	if err != nil {
		return errors.Wrap(err)
	}
	if err := t.opts.Store.Set(string(t.task.Output), data); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// helm renders a helm chart into the output.  The chart is cached per version
// per component and pulled at most once guarded by a filesystem lock.
func (t *taskRunner) helm(ctx context.Context) error {
	h := t.task.Helm
	// Cache the chart by version to pull new versions. (#273)
	cacheDir := filepath.Join(t.opts.AbsLeaf(), "vendor", h.Chart.Version)
	cachePath := filepath.Join(cacheDir, filepath.Base(h.Chart.Name))

	log := logger.FromContext(ctx)

	username := h.Chart.Repository.Auth.Username.Value
	if username == "" {
		username = os.Getenv(h.Chart.Repository.Auth.Username.FromEnv)
	}
	password := h.Chart.Repository.Auth.Password.Value
	if password == "" {
		password = os.Getenv(h.Chart.Repository.Auth.Password.FromEnv)
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		err := func() error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			defer cancel()
			return onceWithLock(log, ctx, cachePath, func() error {
				return errors.Wrap(helm.PullChart(
					ctx,
					cli.New(),
					h.Chart.Name,
					h.Chart.Version,
					h.Chart.Repository.URL,
					cacheDir,
					username,
					password,
				))
			})
		}()
		if err != nil {
			return errors.Format("could not cache chart: %w", err)
		}
	}

	// Write value files
	tempDir, err := os.MkdirTemp("", "holos.helm")
	if err != nil {
		return errors.Format("could not make temp dir: %w", err)
	}
	defer util.Remove(ctx, tempDir)

	// valueFiles represents the ordered list of value files to pass to helm
	// template -f
	var valueFiles []string

	// valueFiles for the use case of migration from helm value hierarchies.
	for _, valueFile := range h.ValueFiles {
		var data []byte
		switch valueFile.Kind {
		case "Values":
			if data, err = yaml.Marshal(valueFile.Values); err != nil {
				return errors.Format("could not marshal value file %s: %w", valueFile.Name, err)
			}
		default:
			return errors.Format("could not marshal value file %s: unknown kind %s", valueFile.Name, valueFile.Kind)
		}

		valuesPath := filepath.Join(tempDir, valueFile.Name)
		if err := os.WriteFile(valuesPath, data, 0666); err != nil {
			return errors.Wrap(fmt.Errorf("could not write value file %s: %w", valueFile.Name, err))
		}
		log.DebugContext(ctx, fmt.Sprintf("wrote: %s", valuesPath))
		valueFiles = append(valueFiles, valuesPath)
	}

	// The final values files
	data, err := yaml.Marshal(h.Values)
	if err != nil {
		return errors.Format("could not marshal values: %w", err)
	}

	valuesPath := filepath.Join(tempDir, "values.yaml")
	if err := os.WriteFile(valuesPath, data, 0666); err != nil {
		return errors.Wrap(fmt.Errorf("could not write values: %w", err))
	}
	log.DebugContext(ctx, fmt.Sprintf("wrote: %s", valuesPath))
	valueFiles = append(valueFiles, valuesPath)

	// Run charts
	args := []string{"template"}
	if !h.EnableHooks {
		args = append(args, "--no-hooks")
	}
	for _, apiVersion := range h.APIVersions {
		args = append(args, "--api-versions", apiVersion)
	}
	if kubeVersion := h.KubeVersion; kubeVersion != "" {
		args = append(args, "--kube-version", kubeVersion)
	}
	args = append(args, "--include-crds")
	for _, valueFilePath := range valueFiles {
		args = append(args, "--values", valueFilePath)
	}
	args = append(args,
		"--namespace", h.Namespace,
		"--kubeconfig", "/dev/null",
		"--version", h.Chart.Version,
		h.Chart.Release,
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
	if err := t.opts.Store.Set(string(t.task.Output), helmOut.Stdout.Bytes()); err != nil {
		return errors.Format("could not store helm output: %w", err)
	}
	log.Debug("set artifact: " + string(t.task.Output))

	return nil
}

// kustomize executes the kustomize command in an isolated temporary directory
// and captures standard output.
func (t *taskRunner) kustomize(ctx context.Context) error {
	store := t.opts.Store
	msg := fmt.Sprintf("could not transform %s for %s", t.task.Output, t.id())

	// Unlike other tasks, kustomize operates in a dedicated temporary directory.
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)

	// Write the kustomization
	data, err := yaml.Marshal(t.task.Kustomize.Kustomization)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	path := filepath.Join(tempDir, "kustomization.yaml")
	if err := os.WriteFile(path, data, 0666); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	// Write additional files, e.g. patch files.
	for name, content := range t.task.Kustomize.Files {
		path := filepath.Join(tempDir, string(name))
		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		if err := os.WriteFile(path, []byte(content), 0666); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
	}

	// Write the inputs
	for _, input := range t.task.Inputs {
		if err := store.Save(tempDir, string(input)); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
	}

	// Execute kustomize
	r, err := util.RunCmdW(ctx, t.opts.Stderr, "kubectl", "kustomize", tempDir)
	if err != nil {
		return errors.Format("%s: could not run kustomize: %w", msg, err)
	}

	// Store the artifact
	if err := store.Set(string(t.task.Output), r.Stdout.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	return nil
}

// join concatenates the inputs into the output with a separator.
func (t *taskRunner) join() error {
	store := t.opts.Store
	s := make([][]byte, 0, len(t.task.Inputs))
	for _, input := range t.task.Inputs {
		if data, ok := store.Get(string(input)); ok {
			s = append(s, data)
		} else {
			return errors.Format("missing input %s", input)
		}
	}
	tempDir, err := t.tempDir()
	if err != nil {
		return errors.Wrap(err)
	}
	// Join the inputs
	data := bytes.Join(s, []byte(t.task.Join.Separator))
	// Save the output to the filesystem.
	outPath := filepath.Join(tempDir, string(t.task.Output))
	if err := os.MkdirAll(filepath.Dir(outPath), 0o777); err != nil {
		return errors.Wrap(err)
	}
	if err := os.WriteFile(outPath, data, 0o666); err != nil {
		return errors.Wrap(err)
	}
	// Store the output in the artifact map.
	if err := store.Load(tempDir, string(t.task.Output)); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// command executes a user defined command with the working directory set to
// the platform root.  Declared inputs are materialized in the build temp dir
// before the command runs.  A command with an output generates or transforms;
// a command with only inputs validates.
func (t *taskRunner) command(ctx context.Context) error {
	store := t.opts.Store

	tempDir, err := t.tempDir()
	if err != nil {
		return errors.Wrap(err)
	}

	args := t.task.Command.Args
	if len(args) < 1 {
		return errors.Format("command args length must be at least 1")
	}

	// Write the inputs.  Materialization goes through sharedSave because
	// concurrent command tasks share the build temp directory and may consume
	// the same input path.
	for _, input := range t.task.Inputs {
		if err := t.sharedSave(tempDir, string(input)); err != nil {
			return errors.Wrap(err)
		}
	}

	// Set the command working directory to the platform root and wire stdin
	// to the named input.
	preRun := func(c *exec.Cmd) error {
		c.Dir = t.opts.Root()
		if stdin := string(t.task.Command.Stdin); stdin != "" {
			data, ok := store.Get(stdin)
			if !ok {
				return errors.Format("stdin input %s not found in the artifact store", stdin)
			}
			c.Stdin = bytes.NewReader(data)
		}
		return nil
	}
	r, err := util.RunCmdFunc(ctx, t.opts.Stderr, args[0], args[1:], preRun)
	if err != nil {
		return errors.Wrap(err)
	}

	// A command with no output validates; there is nothing to store.
	output := string(t.task.Output)
	if output == "" {
		return nil
	}

	// Save the output.
	outPath := filepath.Join(tempDir, output)
	if t.task.Command.IsStdoutOutput {
		if err := os.MkdirAll(filepath.Dir(outPath), 0o777); err != nil {
			return errors.Wrap(err)
		}
		if err := os.WriteFile(outPath, r.Stdout.Bytes(), 0o666); err != nil {
			return errors.Wrap(err)
		}
	}
	if err := store.Load(tempDir, output); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// artifact writes the single input from the artifact store to the final
// artifact path relative to the write-to directory (schema.md D2).
func (t *taskRunner) artifact(ctx context.Context) error {
	store := t.opts.Store
	input := string(t.task.Inputs[0])
	path := string(t.task.Artifact.Path)
	if path == "" {
		path = input
	}

	log := logger.FromContext(ctx)

	if path == input {
		if err := store.Save(t.opts.AbsWriteTo(), path); err != nil {
			return errors.Wrap(err)
		}
		log.DebugContext(ctx, fmt.Sprintf("wrote %s", filepath.Join(t.opts.AbsWriteTo(), path)))
		return nil
	}

	fullPath := filepath.Join(t.opts.AbsWriteTo(), path)

	// A single file input renames to the artifact path.
	if data, ok := store.Get(input); ok {
		if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
			return errors.Wrap(err)
		}
		if err := os.WriteFile(fullPath, data, 0666); err != nil {
			return errors.Wrap(err)
		}
		log.DebugContext(ctx, fmt.Sprintf("wrote %s", fullPath))
		return nil
	}

	// A directory input stages to a temp dir then copies to the artifact path.
	stageDir, err := os.MkdirTemp("", "holos.artifact")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, stageDir)
	if err := store.Save(stageDir, input); err != nil {
		return errors.Wrap(err)
	}
	srcDir := filepath.Join(stageDir, input)
	if _, err := os.Stat(srcDir); err != nil {
		return errors.Format("missing input %s: %w", input, err)
	}
	if err := os.MkdirAll(fullPath, 0777); err != nil {
		return errors.Wrap(err)
	}
	if err := os.CopyFS(fullPath, os.DirFS(srcDir)); err != nil {
		return errors.Wrap(err)
	}
	log.DebugContext(ctx, fmt.Sprintf("wrote %s", fullPath))
	return nil
}

func marshal(list []core.Resource) (buf bytes.Buffer, err error) {
	encoder := yaml.NewEncoder(&buf)
	defer encoder.Close()
	for _, item := range list {
		if err = encoder.Encode(item); err != nil {
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

// BuildContext represents a core BuildContext with version specific helper
// methods.
type BuildContext struct {
	core.BuildContext
}

// Tags returns the cue tags injecting the build context into the TaskSet.
func (bc BuildContext) Tags() ([]string, error) {
	tags := make([]string, 0, 1)
	data, err := json.Marshal(bc.BuildContext)
	if err != nil {
		return tags, errors.Format("could not marshall build context to json: %w", err)
	}
	tags = append(tags, fmt.Sprintf("%s=%s", core.BuildContextTag, string(data)))
	return tags, nil
}

// NewBuildContext returns a new BuildContext
func NewBuildContext(opts holos.BuildOpts) (*BuildContext, error) {
	root := opts.Root()
	if !filepath.IsAbs(root) {
		return nil, errors.Format("not an absolute path: %s", root)
	}
	holosExecutable, err := util.Executable()
	if err != nil {
		return nil, errors.Format("could not get holos executable path: %w", err)
	}
	bc := &BuildContext{
		BuildContext: core.BuildContext{
			TempDir:         opts.TempDir(),
			RootDir:         root,
			LeafDir:         opts.Leaf(),
			HolosExecutable: holosExecutable,
		},
	}
	return bc, nil
}
