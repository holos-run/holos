package v1beta1

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"cuelang.org/go/cue/cuecontext"
	core "github.com/holos-run/holos/api/core/v1beta1"
	"github.com/holos-run/holos/internal/holos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestTaskSet returns a TaskSet over hand-built tasks with a temp platform
// root containing an empty component directory at components/test.
func newTestTaskSet(t *testing.T, tasks map[string]core.Task) *TaskSet {
	t.Helper()
	root := t.TempDir()
	leaf := filepath.Join("components", "test")
	require.NoError(t, os.MkdirAll(filepath.Join(root, leaf), 0o777))
	opts := holos.NewBuildOpts(root, leaf, "deploy", t.TempDir())
	return &TaskSet{
		TaskSet: core.TaskSet{
			APIVersion: "v1beta1",
			Kind:       "TaskSet",
			Metadata:   core.Metadata{Name: "test"},
			Spec:       core.TaskSetSpec{Tasks: tasks},
		},
		Opts: opts,
	}
}

// resourcesTask returns a Resources task producing output with one ConfigMap
// named name.
func resourcesTask(name string, output string) core.Task {
	return core.Task{
		Kind:   "Resources",
		Output: core.FileOrDirectoryPath(output),
		Resources: core.Resources{
			"ConfigMap": {
				core.InternalLabel(name): core.Resource{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata":   map[string]any{"name": name},
				},
			},
		},
	}
}

// recordOrder wires a runHook recording the order tasks start in while still
// executing the real task bodies.
func recordOrder(b *TaskSet) *[]string {
	var mu sync.Mutex
	order := &[]string{}
	b.runHook = func(ctx context.Context, name string, run func(context.Context) error) error {
		mu.Lock()
		*order = append(*order, name)
		mu.Unlock()
		return run(ctx)
	}
	return order
}

func indexOf(t *testing.T, order []string, name string) int {
	t.Helper()
	for idx, val := range order {
		if val == name {
			return idx
		}
	}
	t.Fatalf("task %s not found in order %v", name, order)
	return -1
}

func TestBuildLinearChain(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": resourcesTask("a", "resources.gen.yaml"),
		"combine": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"resources.gen.yaml"},
			Output: "combined.gen.yaml",
		},
		"deploy": {
			Kind:   "Artifact",
			Inputs: []core.FileOrDirectoryPath{"combined.gen.yaml"},
		},
	})
	order := recordOrder(b)

	require.NoError(t, b.Build(t.Context()))

	assert.Less(t, indexOf(t, *order, "gen"), indexOf(t, *order, "combine"))
	assert.Less(t, indexOf(t, *order, "combine"), indexOf(t, *order, "deploy"))

	data, err := os.ReadFile(filepath.Join(b.Opts.AbsWriteTo(), "combined.gen.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "kind: ConfigMap")
}

func TestBuildDiamondDependency(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": resourcesTask("a", "a.gen.yaml"),
		"left": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
			Output: "left.gen.yaml",
		},
		"right": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
			Output: "right.gen.yaml",
		},
		"merge": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"left.gen.yaml", "right.gen.yaml"},
			Join:   core.Join{Separator: "---\n"},
			Output: "merge.gen.yaml",
		},
		"deploy": {
			Kind:     "Artifact",
			Inputs:   []core.FileOrDirectoryPath{"merge.gen.yaml"},
			Artifact: core.Artifact{Path: "components/test/merge.gen.yaml"},
		},
	})
	order := recordOrder(b)

	require.NoError(t, b.Build(t.Context()))

	gen := indexOf(t, *order, "gen")
	merge := indexOf(t, *order, "merge")
	assert.Less(t, gen, indexOf(t, *order, "left"))
	assert.Less(t, gen, indexOf(t, *order, "right"))
	assert.Less(t, indexOf(t, *order, "left"), merge)
	assert.Less(t, indexOf(t, *order, "right"), merge)
	assert.Less(t, merge, indexOf(t, *order, "deploy"))

	data, err := os.ReadFile(filepath.Join(b.Opts.AbsWriteTo(), "components", "test", "merge.gen.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "---\n")
}

func TestBuildCycleError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"alfa": {
			Kind:    "Command",
			Inputs:  []core.FileOrDirectoryPath{"bravo.gen.yaml"},
			Output:  "alfa.gen.yaml",
			Command: core.Command{Args: []string{"true"}},
		},
		"bravo": {
			Kind:    "Command",
			Inputs:  []core.FileOrDirectoryPath{"alfa.gen.yaml"},
			Output:  "bravo.gen.yaml",
			Command: core.Command{Args: []string{"true"}},
		},
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "cycle detected")
	assert.ErrorContains(t, err, "alfa")
	assert.ErrorContains(t, err, "bravo")
}

func TestBuildUnknownInputError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"combine": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"nope.gen.yaml"},
			Output: "combined.gen.yaml",
		},
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "task combine")
	assert.ErrorContains(t, err, "nope.gen.yaml")
	assert.ErrorContains(t, err, "matches no task output")
}

func TestBuildUnknownDependsOnError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": func() core.Task {
			task := resourcesTask("a", "a.gen.yaml")
			task.DependsOn = map[string]core.Dependency{"missing": {}}
			return task
		}(),
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "task gen")
	assert.ErrorContains(t, err, "missing")
	assert.ErrorContains(t, err, "no such task")
}

func TestBuildCanonicalDependsOnError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": func() core.Task {
			task := resourcesTask("a", "a.gen.yaml")
			task.DependsOn = map[string]core.Dependency{"components/vault:helm": {}}
			return task
		}(),
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "canonical task ids are not supported")
}

func TestBuildDuplicateOutputError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"alfa":  resourcesTask("a", "same.gen.yaml"),
		"bravo": resourcesTask("b", "same.gen.yaml"),
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "duplicate output same.gen.yaml")
	assert.ErrorContains(t, err, "alfa")
	assert.ErrorContains(t, err, "bravo")
}

func TestBuildDuplicateArtifactPathError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen-a": resourcesTask("a", "a.gen.yaml"),
		"gen-b": resourcesTask("b", "b.gen.yaml"),
		"deploy-a": {
			Kind:     "Artifact",
			Inputs:   []core.FileOrDirectoryPath{"a.gen.yaml"},
			Artifact: core.Artifact{Path: "same.gen.yaml"},
		},
		"deploy-b": {
			Kind:     "Artifact",
			Inputs:   []core.FileOrDirectoryPath{"b.gen.yaml"},
			Artifact: core.Artifact{Path: "same.gen.yaml"},
		},
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "duplicate artifact path same.gen.yaml")
	assert.ErrorContains(t, err, "deploy-a")
	assert.ErrorContains(t, err, "deploy-b")
}

func TestBuildPathTraversalError(t *testing.T) {
	for _, tc := range []struct {
		name    string
		tasks   map[string]core.Task
		errText string
	}{
		{
			name: "ArtifactPathEscapes",
			tasks: map[string]core.Task{
				"gen": resourcesTask("a", "a.gen.yaml"),
				"deploy": {
					Kind:     "Artifact",
					Inputs:   []core.FileOrDirectoryPath{"a.gen.yaml"},
					Artifact: core.Artifact{Path: "../escape.yaml"},
				},
			},
			errText: "must not traverse outside the write-to directory",
		},
		{
			name: "OutputEscapes",
			tasks: map[string]core.Task{
				"gen": resourcesTask("a", "../a.gen.yaml"),
			},
			errText: "must not traverse outside the build directory",
		},
		{
			name: "OutputAbsolute",
			tasks: map[string]core.Task{
				"gen": resourcesTask("a", "/tmp/a.gen.yaml"),
			},
			errText: "must not traverse outside the build directory",
		},
		{
			name: "InputEscapes",
			tasks: map[string]core.Task{
				"combine": {
					Kind:   "Join",
					Inputs: []core.FileOrDirectoryPath{"../secret.yaml"},
					Output: "combined.gen.yaml",
				},
			},
			errText: "must not traverse outside the build directory",
		},
		{
			name: "FileSourceEscapes",
			tasks: map[string]core.Task{
				"gen": {
					Kind:   "File",
					File:   core.File{Source: "../../etc/passwd"},
					Output: "a.gen.yaml",
				},
			},
			errText: "must not traverse outside the component directory",
		},
		{
			name: "KustomizeFileEscapes",
			tasks: map[string]core.Task{
				"gen": resourcesTask("a", "a.gen.yaml"),
				"transform": {
					Kind:   "Kustomize",
					Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
					Output: "out.gen.yaml",
					Kustomize: core.Kustomize{
						Files: core.FileContentMap{"../patch.yaml": "data"},
					},
				},
			},
			errText: "must not traverse outside the kustomize directory",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := newTestTaskSet(t, tc.tasks)
			err := b.Build(t.Context())
			require.Error(t, err)
			assert.ErrorContains(t, err, tc.errText)
		})
	}
}

func TestBuildSharedCommandInputs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test depends on the cat command")
	}
	// Multiple concurrent command tasks consuming the same input must not
	// race materializing it into the shared build temp directory.
	tasks := map[string]core.Task{
		"gen": resourcesTask("a", "a.gen.yaml"),
	}
	for _, name := range []string{"alfa", "bravo", "charlie", "delta"} {
		tasks[name] = core.Task{
			Kind:   "Command",
			Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
			Output: core.FileOrDirectoryPath(name + ".gen.yaml"),
			Command: core.Command{
				Args:           []string{"cat"},
				Stdin:          "a.gen.yaml",
				IsStdoutOutput: true,
			},
		}
	}
	b := newTestTaskSet(t, tasks)
	b.Opts.Concurrency = 4

	require.NoError(t, b.Build(t.Context()))

	want, ok := b.Opts.Store.Get("a.gen.yaml")
	require.True(t, ok)
	for _, name := range []string{"alfa", "bravo", "charlie", "delta"} {
		have, ok := b.Opts.Store.Get(name + ".gen.yaml")
		require.True(t, ok)
		assert.Equal(t, string(want), string(have))
	}
}

func TestBuildOverlappingArtifactPathError(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen-a": resourcesTask("a", "a.gen.yaml"),
		"gen-b": resourcesTask("b", "b.gen.yaml"),
		"deploy-dir": {
			Kind:     "Artifact",
			Inputs:   []core.FileOrDirectoryPath{"a.gen.yaml"},
			Artifact: core.Artifact{Path: "components/vault"},
		},
		"deploy-file": {
			Kind:     "Artifact",
			Inputs:   []core.FileOrDirectoryPath{"b.gen.yaml"},
			Artifact: core.Artifact{Path: "components/vault/rbac.yaml"},
		},
	})
	err := b.Build(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "overlaps artifact path")
	assert.ErrorContains(t, err, "deploy-dir")
	assert.ErrorContains(t, err, "deploy-file")
}

func TestBuildOverlappingComponentInputs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test depends on the true command")
	}
	// One task consumes the base directory while another consumes a file
	// under it.  The preload must not fail on the write-once store.
	b := newTestTaskSet(t, map[string]core.Task{
		"validate": {
			Kind:    "Command",
			Inputs:  []core.FileOrDirectoryPath{"base"},
			Command: core.Command{Args: []string{"true"}},
		},
		"combine-file": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"base/patch.yaml"},
			Output: "file.gen.yaml",
		},
	})
	dir := filepath.Join(b.Opts.AbsLeaf(), "base")
	require.NoError(t, os.MkdirAll(dir, 0o777))
	want := "greeting: hello\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "patch.yaml"), []byte(want), 0o666))

	require.NoError(t, b.Build(t.Context()))

	data, ok := b.Opts.Store.Get("file.gen.yaml")
	require.True(t, ok)
	assert.Equal(t, want, string(data))
}

func TestBuildArtifactReplacesStaleFiles(t *testing.T) {
	newDirTaskSet := func(t *testing.T, artifactPath string) *TaskSet {
		b := newTestTaskSet(t, map[string]core.Task{
			"deploy": {
				Kind:     "Artifact",
				Inputs:   []core.FileOrDirectoryPath{"manifests"},
				Artifact: core.Artifact{Path: core.FileOrDirectoryPath(artifactPath)},
			},
		})
		dir := filepath.Join(b.Opts.AbsLeaf(), "manifests")
		require.NoError(t, os.MkdirAll(dir, 0o777))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "real.yaml"), []byte("real: true\n"), 0o666))
		return b
	}

	for _, tc := range []struct {
		name string
		path string
		dest string
	}{
		{name: "DefaultPath", path: "", dest: "manifests"},
		{name: "RenamedPath", path: "out", dest: "out"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := newDirTaskSet(t, tc.path)
			// Simulate a stale file from a previous render.
			staleDir := filepath.Join(b.Opts.AbsWriteTo(), tc.dest)
			require.NoError(t, os.MkdirAll(staleDir, 0o777))
			require.NoError(t, os.WriteFile(filepath.Join(staleDir, "stale.yaml"), []byte("stale: true\n"), 0o666))

			require.NoError(t, b.Build(t.Context()))

			_, err := os.Stat(filepath.Join(staleDir, "real.yaml"))
			assert.NoError(t, err, "expected current file to be written")
			_, err = os.Stat(filepath.Join(staleDir, "stale.yaml"))
			assert.True(t, os.IsNotExist(err), "expected stale file to be removed")
		})
	}
}

func TestBuildDisabledNoOp(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": resourcesTask("a", "a.gen.yaml"),
		"deploy": {
			Kind:   "Artifact",
			Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
		},
	})
	b.Spec.Disabled = true
	order := recordOrder(b)

	require.NoError(t, b.Build(t.Context()))

	assert.Empty(t, *order, "expected no tasks to run for a disabled TaskSet")
	_, err := os.Stat(filepath.Join(b.Opts.AbsWriteTo(), "a.gen.yaml"))
	assert.True(t, os.IsNotExist(err), "expected no artifact to be written")
}

func TestBuildDependsOnOrdering(t *testing.T) {
	second := resourcesTask("b", "b.gen.yaml")
	second.DependsOn = map[string]core.Dependency{"first": {}}
	b := newTestTaskSet(t, map[string]core.Task{
		"first":  resourcesTask("a", "a.gen.yaml"),
		"second": second,
	})
	b.Opts.Concurrency = 1
	order := recordOrder(b)

	require.NoError(t, b.Build(t.Context()))

	assert.Less(t, indexOf(t, *order, "first"), indexOf(t, *order, "second"))
}

func TestBuildComponentDirectoryInput(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"combine": {
			Kind:   "Join",
			Inputs: []core.FileOrDirectoryPath{"input.yaml"},
			Output: "combined.gen.yaml",
		},
		"deploy": {
			Kind:   "Artifact",
			Inputs: []core.FileOrDirectoryPath{"combined.gen.yaml"},
		},
	})
	want := "greeting: hello\n"
	path := filepath.Join(b.Opts.AbsLeaf(), "input.yaml")
	require.NoError(t, os.WriteFile(path, []byte(want), 0o666))

	require.NoError(t, b.Build(t.Context()))

	data, err := os.ReadFile(filepath.Join(b.Opts.AbsWriteTo(), "combined.gen.yaml"))
	require.NoError(t, err)
	assert.Equal(t, want, string(data))
}

func TestBuildDeterministicReadyOrder(t *testing.T) {
	tasks := make(map[string]core.Task)
	names := []string{"alfa", "bravo", "charlie", "delta", "echo"}
	for _, name := range names {
		tasks[name] = resourcesTask(name, name+".gen.yaml")
	}
	// Repeat the build to guard against map iteration order masking a bug.
	for i := 0; i < 5; i++ {
		b := newTestTaskSet(t, tasks)
		b.Opts.Concurrency = 1
		order := recordOrder(b)
		require.NoError(t, b.Build(t.Context()))
		assert.Equal(t, names, *order, "expected ready tasks to run in sorted name order")
	}
}

func TestBuildConcurrencyBound(t *testing.T) {
	tasks := make(map[string]core.Task)
	for _, name := range []string{"alfa", "bravo", "charlie", "delta", "echo", "foxtrot"} {
		tasks[name] = resourcesTask(name, name+".gen.yaml")
	}

	track := func(b *TaskSet) *int {
		var mu sync.Mutex
		var current int
		maxSeen := new(int)
		b.runHook = func(ctx context.Context, name string, run func(context.Context) error) error {
			mu.Lock()
			current++
			if current > *maxSeen {
				*maxSeen = current
			}
			mu.Unlock()
			time.Sleep(20 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
			return nil
		}
		return maxSeen
	}

	t.Run("Bounded", func(t *testing.T) {
		b := newTestTaskSet(t, tasks)
		b.Opts.Concurrency = 2
		maxSeen := track(b)
		require.NoError(t, b.Build(t.Context()))
		assert.LessOrEqual(t, *maxSeen, 2, "expected at most 2 tasks in flight")
		assert.GreaterOrEqual(t, *maxSeen, 2, "expected independent tasks to run concurrently")
	})

	t.Run("Sequential", func(t *testing.T) {
		b := newTestTaskSet(t, tasks)
		b.Opts.Concurrency = 1
		maxSeen := track(b)
		require.NoError(t, b.Build(t.Context()))
		assert.Equal(t, 1, *maxSeen, "expected exactly 1 task in flight")
	})

	t.Run("ZeroDefaultsToOne", func(t *testing.T) {
		b := newTestTaskSet(t, tasks)
		b.Opts.Concurrency = 0
		maxSeen := track(b)
		require.NoError(t, b.Build(t.Context()))
		assert.Equal(t, 1, *maxSeen, "expected exactly 1 task in flight")
	})
}

func TestBuildTaskFailureStopsBuild(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": {
			Kind:   "File",
			File:   core.File{Source: "does-not-exist.yaml"},
			Output: "a.gen.yaml",
		},
		"deploy": {
			Kind:   "Artifact",
			Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
		},
	})
	order := recordOrder(b)

	err := b.Build(t.Context())
	require.Error(t, err)
	assert.NotContains(t, *order, "deploy", "expected downstream task to be skipped on failure")
}

func TestBuildCommandStdin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test depends on the cat command")
	}
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": resourcesTask("a", "a.gen.yaml"),
		"copy": {
			Kind:   "Command",
			Inputs: []core.FileOrDirectoryPath{"a.gen.yaml"},
			Output: "copy.gen.yaml",
			Command: core.Command{
				Args:           []string{"cat"},
				Stdin:          "a.gen.yaml",
				IsStdoutOutput: true,
			},
		},
		"deploy": {
			Kind:   "Artifact",
			Inputs: []core.FileOrDirectoryPath{"copy.gen.yaml"},
		},
	})

	require.NoError(t, b.Build(t.Context()))

	want, ok := b.Opts.Store.Get("a.gen.yaml")
	require.True(t, ok)
	data, err := os.ReadFile(filepath.Join(b.Opts.AbsWriteTo(), "copy.gen.yaml"))
	require.NoError(t, err)
	assert.Equal(t, string(want), string(data))
}

func TestBuildCommandValidator(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test depends on the false command")
	}
	// A command with only inputs validates, gating downstream tasks through
	// dependsOn edges (schema.md D2).
	deploy := core.Task{
		Kind:      "Artifact",
		Inputs:    []core.FileOrDirectoryPath{"a.gen.yaml"},
		DependsOn: map[string]core.Dependency{"validate": {}},
	}

	t.Run("Pass", func(t *testing.T) {
		b := newTestTaskSet(t, map[string]core.Task{
			"gen": resourcesTask("a", "a.gen.yaml"),
			"validate": {
				Kind:    "Command",
				Inputs:  []core.FileOrDirectoryPath{"a.gen.yaml"},
				Command: core.Command{Args: []string{"true"}},
			},
			"deploy": deploy,
		})
		require.NoError(t, b.Build(t.Context()))
		_, err := os.Stat(filepath.Join(b.Opts.AbsWriteTo(), "a.gen.yaml"))
		assert.NoError(t, err, "expected artifact written after validation passed")
	})

	t.Run("Fail", func(t *testing.T) {
		b := newTestTaskSet(t, map[string]core.Task{
			"gen": resourcesTask("a", "a.gen.yaml"),
			"validate": {
				Kind:    "Command",
				Inputs:  []core.FileOrDirectoryPath{"a.gen.yaml"},
				Command: core.Command{Args: []string{"false"}},
			},
			"deploy": deploy,
		})
		err := b.Build(t.Context())
		require.Error(t, err)
		_, statErr := os.Stat(filepath.Join(b.Opts.AbsWriteTo(), "a.gen.yaml"))
		assert.True(t, os.IsNotExist(statErr), "expected no artifact written after validation failed")
	})
}

func TestMatchProducers(t *testing.T) {
	producers := map[string]string{
		"a.gen.yaml":       "alfa",
		"out/one.yaml":     "bravo",
		"out/two.yaml":     "charlie",
		"manifests":        "delta",
		"other/three.yaml": "echo",
	}
	outputs := []string{"a.gen.yaml", "manifests", "other/three.yaml", "out/one.yaml", "out/two.yaml"}

	for _, tc := range []struct {
		name  string
		input string
		want  []string
	}{
		{name: "ExactMatch", input: "a.gen.yaml", want: []string{"alfa"}},
		{name: "DirectoryInput", input: "out", want: []string{"bravo", "charlie"}},
		{name: "DirectoryOutput", input: "manifests/deep/file.yaml", want: []string{"delta"}},
		{name: "NoMatch", input: "missing.yaml", want: nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, matchProducers(tc.input, producers, outputs))
		})
	}
}

func TestValidateTask(t *testing.T) {
	for _, tc := range []struct {
		name    string
		task    core.Task
		errText string
	}{
		{
			name:    "InvalidName",
			task:    core.Task{Kind: "Resources", Output: "a.yaml"},
			errText: "invalid task name",
		},
		{
			name:    "UnsupportedKind",
			task:    core.Task{Kind: "Bogus"},
			errText: "unsupported kind Bogus",
		},
		{
			name:    "HelmWithInputs",
			task:    core.Task{Kind: "Helm", Inputs: []core.FileOrDirectoryPath{"a.yaml"}, Output: "b.yaml"},
			errText: "must not declare inputs",
		},
		{
			name:    "ResourcesWithoutOutput",
			task:    core.Task{Kind: "Resources"},
			errText: "requires an output",
		},
		{
			name:    "JoinWithoutInputs",
			task:    core.Task{Kind: "Join", Output: "a.yaml"},
			errText: "requires at least one input",
		},
		{
			name:    "CommandWithoutArgs",
			task:    core.Task{Kind: "Command"},
			errText: "args length must be at least 1",
		},
		{
			name: "CommandStdoutWithoutOutput",
			task: core.Task{
				Kind:    "Command",
				Command: core.Command{Args: []string{"true"}, IsStdoutOutput: true},
			},
			errText: "requires an output when isStdoutOutput",
		},
		{
			name: "CommandStdinNotDeclared",
			task: core.Task{
				Kind:    "Command",
				Command: core.Command{Args: []string{"cat"}, Stdin: "a.yaml"},
			},
			errText: "must be one of the task inputs",
		},
		{
			name:    "ArtifactTwoInputs",
			task:    core.Task{Kind: "Artifact", Inputs: []core.FileOrDirectoryPath{"a.yaml", "b.yaml"}},
			errText: "exactly one input",
		},
		{
			name:    "ArtifactWithOutput",
			task:    core.Task{Kind: "Artifact", Inputs: []core.FileOrDirectoryPath{"a.yaml"}, Output: "b.yaml"},
			errText: "must not declare an output",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := "valid-name"
			if tc.name == "InvalidName" {
				name = "Invalid_Name"
			}
			err := validateTask(name, tc.task)
			require.Error(t, err)
			assert.ErrorContains(t, err, tc.errText)
		})
	}

	t.Run("Valid", func(t *testing.T) {
		err := validateTask("valid-name", core.Task{Kind: "Resources", Output: "a.yaml"})
		assert.NoError(t, err)
	})
}

func TestLoad(t *testing.T) {
	cuectx := cuecontext.New()

	t.Run("Valid", func(t *testing.T) {
		v := cuectx.CompileString(`{
			apiVersion: "v1beta1"
			kind: "TaskSet"
			metadata: name: "test"
			spec: tasks: gen: {
				kind: "Resources"
				output: "a.gen.yaml"
				resources: {}
			}
			buildContext: {
				tempDir: "/tmp/holos"
				rootDir: "/platform"
				leafDir: "components/test"
				holosExecutable: "holos"
			}
		}`)
		var b TaskSet
		require.NoError(t, b.Load(v))
		assert.Equal(t, "test", b.Metadata.Name)
		assert.Equal(t, "Resources", b.Spec.Tasks["gen"].Kind)
	})

	t.Run("NotConcrete", func(t *testing.T) {
		v := cuectx.CompileString(`{
			apiVersion: "v1beta1"
			kind: "TaskSet"
			metadata: name: string
			spec: tasks: {}
			buildContext: {
				tempDir: "/tmp/holos"
				rootDir: "/platform"
				leafDir: "components/test"
				holosExecutable: "holos"
			}
		}`)
		var b TaskSet
		assert.Error(t, b.Load(v))
	})
}

func TestExport(t *testing.T) {
	b := newTestTaskSet(t, map[string]core.Task{
		"gen": resourcesTask("a", "a.gen.yaml"),
	})
	var buf bytes.Buffer
	encoder, err := holos.NewSequentialEncoder("yaml", &buf)
	require.NoError(t, err)
	require.NoError(t, b.Export(0, encoder))
	require.NoError(t, encoder.Close())
	assert.Contains(t, buf.String(), "kind: TaskSet")
	assert.Contains(t, buf.String(), "apiVersion: v1beta1")
}

func TestBuildContextTags(t *testing.T) {
	root := t.TempDir()
	opts := holos.NewBuildOpts(root, "components/test", "deploy", t.TempDir())
	bc, err := NewBuildContext(opts)
	require.NoError(t, err)
	tags, err := bc.Tags()
	require.NoError(t, err)
	require.Len(t, tags, 1)
	assert.Contains(t, tags[0], fmt.Sprintf("%s=", core.BuildContextTag))
	assert.Contains(t, tags[0], root)

	t.Run("RelativeRootError", func(t *testing.T) {
		opts := holos.NewBuildOpts("relative/root", "components/test", "deploy", t.TempDir())
		_, err := NewBuildContext(opts)
		assert.Error(t, err)
	})
}
