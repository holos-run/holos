# v1alpha6 Build Plan Executor

This document outlines the plan to finalize the v1alpha6 build plan executor for `holos render platform`. The goal is to use the compiler pool to produce a complete set of BuildPlans, then combine them into a topologically sorted DAG for efficient concurrent execution.

## Related Issues

- [#435 v1alpha6 Design](https://github.com/holos-run/holos/issues/435)
- [#389 Concurrent Component Compiler](https://github.com/holos-run/holos/issues/389)
- [#452 Fix concurrent helm invocations](https://github.com/holos-run/holos/issues/452)

## Current State Analysis

### What's Implemented

1. **Compiler Pool** (`internal/compile/compile.go`):
   - `compile.Compile()` function spawns a pool of `holos compile` subprocesses
   - Producer-consumer pattern with task queue using channels
   - Each worker maintains a persistent subprocess with stdin/stdout pipes
   - JSON streaming protocol for BuildPlanRequest/BuildPlanResponse

2. **`holos show buildplans`** (`internal/cli/show.go`):
   - Uses `compile.Compile()` to produce BuildPlans concurrently
   - Displays BuildPlans without executing them
   - Validates the compiler pool works correctly

3. **`holos render platform`** (`internal/cli/render/render.go`):
   - Spawns individual `holos render component` subprocesses
   - Each subprocess independently compiles AND executes its BuildPlan
   - Does NOT use the compiler pool

4. **BuildPlan Execution** (`internal/component/v1alpha6/v1alpha6.go`):
   - `BuildPlan.Build()` processes artifacts concurrently
   - Within each artifact: generators run concurrently, transformers sequentially, validators concurrently
   - Worker pool for task execution within a single BuildPlan

### What's Missing

1. **Compiler Pool in `holos render platform`**: The render command still spawns individual subprocesses instead of using the compiler pool.

2. **Unified DAG Execution**: BuildPlans execute independently. No cross-BuildPlan task deduplication or optimization.

3. **Helm Chart Deduplication**: Each component fetches Helm charts independently. Multiple components using the same chart/version fetch redundantly.

4. **Topological Sort**: No dependency analysis or task ordering across BuildPlans.

## Architecture Design

### Phase 1: Compile All BuildPlans

Use the existing compiler pool to compile all BuildPlans before executing any.

```
holos render platform
  │
  ├─> Platform.Build()
  │     └─> Collect all Components
  │
  ├─> compile.Compile(ctx, concurrency, requests)
  │     ├─> Producer: sends BuildPlanRequest per component
  │     ├─> N Workers: each runs holos compile subprocess
  │     └─> Returns: []BuildPlanResponse with all BuildPlans
  │
  └─> Phase 2: Execute BuildPlans
```

### Phase 2: Build Unified Task DAG

Transform all BuildPlans into a single DAG of executable tasks.

```go
// Task represents a single executable unit of work
type Task struct {
    ID       string              // Unique identifier
    Kind     string              // helm-fetch, generator, transformer, validator, write-artifact
    BuildPlan string             // Which BuildPlan this task belongs to
    Artifact  string             // Which Artifact within the BuildPlan
    Inputs   []string            // Task IDs this depends on
    Outputs  []string            // Output keys this produces
    Run      func(ctx) error     // Execution function
}

// TaskDAG represents the complete dependency graph
type TaskDAG struct {
    Tasks     map[string]*Task   // All tasks by ID
    Adjacency map[string][]string // task -> tasks that depend on it
    InDegree  map[string]int     // task -> number of dependencies
}
```

#### Task Types

1. **helm-fetch**: Download a Helm chart to local cache
   - Key: `helm-fetch:{repo}:{chart}:{version}`
   - Deduplicated across all BuildPlans

2. **generator**: Execute a Generator (Resources, Helm template, File, Command)
   - Key: `{buildplan}:{artifact}:gen:{output}`
   - helm-fetch tasks depend on generators that use Helm

3. **transformer**: Execute a Transformer (Kustomize, Join, Command)
   - Key: `{buildplan}:{artifact}:transform:{output}`
   - Sequential within artifact, depends on generators or prior transformers

4. **validator**: Execute a Validator (Command)
   - Key: `{buildplan}:{artifact}:validate:{idx}`
   - Depends on final transformer output

5. **write-artifact**: Write artifact to filesystem
   - Key: `{buildplan}:{artifact}:write`
   - Depends on validators completing

### Phase 3: Topological Sort and Execution

Execute tasks using Kahn's algorithm for topological ordering with concurrent execution.

```go
func ExecuteDAG(ctx context.Context, dag *TaskDAG, workers int) error {
    g, ctx := errgroup.WithContext(ctx)
    ready := make(chan *Task, len(dag.Tasks))
    completed := make(chan string, len(dag.Tasks))

    // Initialize with tasks that have no dependencies
    for id, task := range dag.Tasks {
        if dag.InDegree[id] == 0 {
            ready <- task
        }
    }

    // Coordinator goroutine
    g.Go(func() error {
        remaining := len(dag.Tasks)
        for remaining > 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case id := <-completed:
                remaining--
                // Decrement in-degree of dependent tasks
                for _, depID := range dag.Adjacency[id] {
                    dag.InDegree[depID]--
                    if dag.InDegree[depID] == 0 {
                        ready <- dag.Tasks[depID]
                    }
                }
            }
        }
        close(ready)
        return nil
    })

    // Worker goroutines
    for i := 0; i < workers; i++ {
        g.Go(func() error {
            for task := range ready {
                if err := task.Run(ctx); err != nil {
                    return err
                }
                completed <- task.ID
            }
            return nil
        })
    }

    return g.Wait()
}
```

### Helm Chart Deduplication

Extract all Helm chart references, deduplicate, and create fetch tasks.

```go
type ChartKey struct {
    Repository string
    Name       string
    Version    string
}

func ExtractHelmCharts(plans []BuildPlan) map[ChartKey][]TaskRef {
    charts := make(map[ChartKey][]TaskRef)
    for _, plan := range plans {
        for _, artifact := range plan.Spec.Artifacts {
            for _, gen := range artifact.Generators {
                if gen.Kind == "Helm" {
                    key := ChartKey{
                        Repository: gen.Helm.Chart.Repository.URL,
                        Name:       gen.Helm.Chart.Name,
                        Version:    gen.Helm.Chart.Version,
                    }
                    charts[key] = append(charts[key], TaskRef{
                        BuildPlan: plan.Metadata.Name,
                        Artifact:  artifact.Artifact,
                        Generator: gen.Output,
                    })
                }
            }
        }
    }
    return charts
}
```

## Implementation Plan

### Step 1: Wire Compiler Pool to Render Platform

Modify `internal/cli/render/render.go` to:

1. Collect all components from the Platform
2. Build `[]BuildPlanRequest` for each component
3. Call `compile.Compile()` to get all BuildPlans
4. Execute BuildPlans (initially using existing per-BuildPlan execution)

**Files to modify:**
- `internal/cli/render/render.go`

**Validation:**
- `holos render platform` produces identical output to before
- Compilation is faster due to concurrent pool

### Step 2: Create Task DAG Package

Create new package `internal/dag/` with:

1. Task and TaskDAG types
2. DAG construction from BuildPlans
3. Topological sort execution
4. Shared artifact store across BuildPlans

**Files to create:**
- `internal/dag/dag.go` - Core types and construction
- `internal/dag/execute.go` - DAG execution with worker pool
- `internal/dag/dag_test.go` - Unit tests

### Step 3: Implement Helm Chart Deduplication

1. Extract all Helm chart references before execution
2. Create single fetch task per unique chart
3. Helm generator tasks depend on corresponding fetch task
4. Use file system locking for concurrent safety

**Files to modify:**
- `internal/dag/dag.go` - Add chart extraction and fetch task creation
- `internal/component/v1alpha6/v1alpha6.go` - Refactor helm fetch logic

### Step 4: Integrate DAG Execution into Render Platform

Replace per-BuildPlan execution with unified DAG execution.

**Files to modify:**
- `internal/cli/render/render.go`
- `internal/platform/platform.go`

### Step 5: Update Build Plan Execution to Use Shared Store

Modify BuildPlan execution to:

1. Accept a shared artifact store
2. Support external Helm chart cache
3. Report task completion to coordinator

**Files to modify:**
- `internal/component/v1alpha6/v1alpha6.go`
- `internal/artifact/artifact.go`

## Testing Strategy

### Unit Tests

1. DAG construction from BuildPlans
2. Topological sort correctness
3. Cycle detection
4. Helm chart deduplication

### Integration Tests

1. Render platform with multiple components sharing Helm charts
2. Verify deduplication reduces Helm fetches
3. Verify output matches non-DAG execution
4. Use `holos compare buildplans` to validate equivalence

### Performance Tests

1. Measure compilation time improvement
2. Measure execution time with chart deduplication
3. Benchmark with varying concurrency levels

## Migration Path

1. **v0.106.0**: Wire compiler pool to render platform (Step 1)
2. **v0.107.0**: Introduce DAG package, optional via flag (Steps 2-3)
3. **v0.108.0**: Make DAG execution default, deprecate old path (Steps 4-5)

## Future Considerations

### TaskSet Schema (v1alpha6 Design Goal)

The DAG implementation provides a foundation for the TaskSet schema:

```cue
TaskSet: {
    tasks: [string]: Task
}

Task: {
    kind: "HelmFetch" | "Generator" | "Transformer" | "Validator" | "Write"
    inputs: [...string]  // Task IDs
    outputs: [...string] // Output keys
    // Kind-specific fields
}
```

### Cross-Component Dependencies

Future enhancement to support:

```cue
Component: {
    dependsOn: [...string]  // Other component names
}
```

### Caching

Consider caching:

1. Compiled BuildPlans (CUE evaluation cache)
2. Generator outputs (content-addressed)
3. Transformer outputs (input-hash based)

## Appendix: Current Code Flow

### `holos show buildplans` (Uses Compiler Pool)

```
show.go:84  showBuildPlans.Run()
  │
  ├─> p.Select() - Get components matching selectors
  │
  ├─> Build []BuildPlanRequest from components
  │
  ├─> compile.Compile(ctx, concurrency, reqs)
  │     compile.go:155
  │     ├─> Producer goroutine sends requests to channel
  │     ├─> N Consumer goroutines run holos compile
  │     └─> Returns []BuildPlanResponse
  │
  └─> Encode and output BuildPlans
```

### `holos render platform` (Current: Per-Component Subprocess)

```
render.go:49  renderPlatform.Run()
  │
  ├─> platform.Platform.Build()
  │     platform.go
  │     └─> For each component (concurrent, limited):
  │           └─> PerComponentFunc()
  │
  └─> PerComponentFunc() spawns subprocess:
        holos render component [path]
        │
        └─> component.go
              ├─> Compile CUE to BuildPlan
              └─> BuildPlan.Build() - Execute artifact pipeline
```

### Proposed: `holos render platform` (With Compiler Pool + DAG)

```
render.go  renderPlatform.Run()
  │
  ├─> Collect all components from Platform
  │
  ├─> compile.Compile(ctx, concurrency, reqs)
  │     └─> Returns []BuildPlanResponse
  │
  ├─> dag.Build(buildPlans)
  │     ├─> Extract and deduplicate Helm charts
  │     ├─> Create fetch tasks
  │     ├─> Create generator tasks (depend on fetch)
  │     ├─> Create transformer tasks (depend on generators)
  │     ├─> Create validator tasks (depend on transformers)
  │     └─> Create write tasks (depend on validators)
  │
  └─> dag.Execute(ctx, workers)
        ├─> Topological sort execution
        ├─> Fixed worker pool with errgroup
        └─> Write artifacts to filesystem
```
