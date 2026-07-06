# Schema: TaskSet, Task, and Command

This chapter specifies the v1beta1 core schema: the `TaskSet` resource that
replaces the deprecated BuildPlan, the `Task` as the single unit of work in a
data transformation pipeline (subsuming the v1alpha6 Generator, Transformer,
and Validator concepts), and the first-class `Command` task kind for invoking
external tools.  It is the design authority for Phase 1: the struct
definitions below are transcribed into `api/core/v1beta1/types.go`, and the
numbered design decisions
([D1](#d1-edge-derivation)–[D5](#d5-open-and-closed-structs)) are cited by
later phases as `schema.md#d1-edge-derivation` and so on.  The properties
fixed here — struct-keyed tasks, explicit dependency edges, component-path
namespacing — are what let [rendering.md](rendering.md) compose every
component's TaskSet into one platform-wide DAG.

## Before: the v1alpha6 shape being replaced

A v1alpha6 BuildPlan holds a **list** of artifacts, and each artifact holds
three **phase-ordered lists** of wrapper kinds.  From
[`api/core/v1alpha6/types.go`](../../../api/core/v1alpha6/types.go)
(lines 131–174):

```go
// BuildPlanSpec represents the specification of the [BuildPlan].
type BuildPlanSpec struct {
	// Artifacts represents the artifacts for holos to build.
	Artifacts []Artifact `json:"artifacts" yaml:"artifacts"`
	// Disabled causes the holos render platform command to skip the BuildPlan.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

type Artifact struct {
	Artifact     FileOrDirectoryPath `json:"artifact,omitempty" yaml:"artifact,omitempty"`
	Generators   []Generator         `json:"generators,omitempty" yaml:"generators,omitempty"`
	Transformers []Transformer       `json:"transformers,omitempty" yaml:"transformers,omitempty"`
	Validators   []Validator         `json:"validators,omitempty" yaml:"validators,omitempty"`
	Skip         bool                `json:"skip,omitempty" yaml:"skip,omitempty"`
}
```

The execution semantics are implicit in the shape: generators run
concurrently, transformers run sequentially in list order, validators run
concurrently after the final transformer, and then holos writes the
`artifact:` path.  The `Generator`, `Transformer`, and `Validator` wrapper
kinds exist only to encode those phases — each wraps the same underlying
config kinds (`Helm`, `Kustomize`, `Command`, …) with slightly different
input/output fields.

`Command` today is a leaf config shared by the three wrapper kinds, not a
task kind itself (`api/core/v1alpha6/types.go` lines 371–383):

```go
// Command represents a [BuildPlan] task implemented by executing an user
// defined system command.  A task is defined as a [Generator], [Transformer],
// or [Validator].  Commands are executed with the working directory set to the
// platform root.
type Command struct {
	// DisplayName of the command.  The basename of args[0] is used if empty.
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	// Args represents the argument vector passed to the os. to execute the
	// command.
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
	// IsStdoutOutput captures the command stdout as the task output if true.
	IsStdoutOutput bool `json:"isStdoutOutput,omitempty" yaml:"isStdoutOutput,omitempty"`
}
```

Three properties of this shape motivate the redesign:

1. **Lists do not unify.**  CUE list unification is positional, so two files
   cannot both append a transformer to the same artifact (design-inputs
   item 2 in the [README](README.md#design-inputs)).  A mixin that wants to
   add a validator to another component's pipeline has no place to put it.
2. **The DAG is implicit and fixed.**  Phase ordering
   (generate → transform → validate → write) is baked into the wrapper
   kinds; any other pipeline shape cannot be expressed, and cross-component
   edges are impossible, blocking the platform-wide DAG (items 4–5).
3. **Three kinds for one concept.**  Generator, Transformer, and Validator
   differ only in phase and in which input/output fields they carry.  The
   duplication triples the schema surface and forces `Command` to be a leaf
   config rather than a task.

## After: the v1beta1 core schema

v1beta1 collapses the three wrapper kinds into one `Task`, keys tasks by
name in a struct, and makes the DAG explicit.  The following definitions are
normative; Phase 1 (HOL-1493) transcribes them into
`api/core/v1beta1/types.go`.

```go
// TaskSet replaces BuildPlan.  A component produces one TaskSet; holos merges
// all component TaskSets into one platform-wide DAG.
type TaskSet struct {
	// APIVersion represents the versioned schema of the resource.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"\"v1beta1\""`
	// Kind represents the type of the resource.
	Kind string `json:"kind" yaml:"kind" cue:"\"TaskSet\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	// Spec specifies the desired state of the resource.
	Spec TaskSetSpec `json:"spec" yaml:"spec"`
	// BuildContext represents values injected by holos just before evaluating
	// a TaskSet, for example the tempDir used for the build.
	BuildContext BuildContext `json:"buildContext" yaml:"buildContext"`
}

// TaskSetSpec represents the specification of the [TaskSet].
type TaskSetSpec struct {
	// Tasks keyed by name — structs, not lists, so TaskSets compose by CUE
	// unification (design item 2; author-schema NameLabel idiom).
	Tasks map[string]Task `json:"tasks" yaml:"tasks"`
	// Disabled causes the holos render platform command to skip the TaskSet.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// Task unifies the v1alpha6 Generator, Transformer, and Validator concepts.
// A task declares the artifact-store paths it consumes and produces; the
// executor derives DAG edges from those declarations (see D1).
type Task struct {
	// Kind discriminates the task behavior.
	Kind string `json:"kind" yaml:"kind" cue:"\"Resources\" | \"Helm\" | \"File\" | \"Kustomize\" | \"Join\" | \"Command\" | \"Artifact\""`
	// DependsOn declares tasks that must complete before this task runs,
	// keyed by task name or canonical ID (see D3) — a struct, not a list, so
	// mixins compose ordering edges by unification.  Use for ordering
	// constraints with no data flow; data-flow edges are derived from Inputs
	// and Output (see D1).
	DependsOn map[string]Dependency `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`
	// Inputs are artifact-store paths consumed by the task.
	Inputs []FileOrDirectoryPath `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	// Output is the artifact-store path produced by the task.  Output values
	// are write-once: it is an error for two tasks to declare the same
	// Output within one platform render.
	Output FileOrDirectoryPath `json:"output,omitempty" yaml:"output,omitempty"`

	// Exactly one of the following, matching Kind:
	Resources Resources `json:"resources,omitempty" yaml:"resources,omitempty"`
	Helm      Helm      `json:"helm,omitempty" yaml:"helm,omitempty"`
	File      File      `json:"file,omitempty" yaml:"file,omitempty"`
	Kustomize Kustomize `json:"kustomize,omitempty" yaml:"kustomize,omitempty"`
	Join      Join      `json:"join,omitempty" yaml:"join,omitempty"`
	Command   Command   `json:"command,omitempty" yaml:"command,omitempty"`
	Artifact  Artifact  `json:"artifact,omitempty" yaml:"artifact,omitempty"`
}

// Dependency represents one explicit ordering edge declared in
// [Task.DependsOn].  It is deliberately empty — the edge is the struct key —
// so future fields (for example an optional edge) may be added without
// breaking composition.
type Dependency struct{}

// Command is a first-class Task kind in v1beta1.  Commands execute with the
// working directory set to the platform root.
type Command struct {
	// DisplayName of the command.  The basename of args[0] is used if empty.
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	// Args represents the argument vector passed to the os to execute the
	// command.
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
	// Stdin names a task input wired to the command's standard input.  Must
	// be one of the task's declared Inputs.
	Stdin FileOrDirectoryPath `json:"stdin,omitempty" yaml:"stdin,omitempty"`
	// IsStdoutOutput captures the command stdout as the task Output if true.
	IsStdoutOutput bool `json:"isStdoutOutput,omitempty" yaml:"isStdoutOutput,omitempty"`
}

// Artifact is the sink Task kind.  It writes its single input from the
// artifact store to the final artifact path once every task it depends on
// has completed successfully (see D2).
type Artifact struct {
	// Path represents the final artifact path relative to the write-to
	// directory, e.g. deploy/.  Defaults to the task's single input path
	// when empty.
	Path FileOrDirectoryPath `json:"path,omitempty" yaml:"path,omitempty"`
}
```

The `Resources`, `Helm`, `File`, `Kustomize`, and `Join` config structs carry
over from v1alpha6 with their existing fields; only their wrapper changes
(from Generator/Transformer to Task).  `Metadata` and `FileOrDirectoryPath`
carry over unchanged.

### Task kinds

The full v1beta1 task kind list, with the v1alpha6 concept each subsumes
and the normative I/O cardinality of each kind:

| Kind | Behavior | `inputs` | `output` | v1alpha6 ancestor |
| -- | -- | -- | -- | -- |
| `Resources` | export Kubernetes resources defined in CUE | none | required | Generator |
| `Helm` | render a Helm chart | none | required | Generator |
| `File` | read a file from the component directory | none | required | Generator |
| `Kustomize` | patch and transform prior outputs | one or more | required | Transformer |
| `Join` | concatenate prior outputs | one or more | required | Transformer |
| `Command` | execute a user-defined command | zero or more | optional; required when `isStdoutOutput` | Generator, Transformer, and Validator leaf config |
| `Artifact` | write the final artifact (sink; see [D2](#d2-artifact-writing)) | exactly one | none | the implicit `artifact:` write |

The `inputs` and `output` columns are normative requiredness and cardinality
constraints, enforced by the same per-kind CUE guards that enforce the
config field ([D5](#d5-open-and-closed-structs)) and revalidated by the
executor so schema and runtime cannot diverge: a `Helm` task declaring
`inputs`, a `Kustomize` task without them, a `Command` task setting
`isStdoutOutput: true` without an `output`, or an `Artifact` task with two
inputs all fail evaluation.

There is no `Validator` kind: a validator is a `Command` task that declares
`inputs` and no `output`, gating downstream tasks through `dependsOn` edges
([D1](#d1-edge-derivation), [D2](#d2-artifact-writing)).

### Command as a first-class task kind

In v1alpha6, running a tool requires choosing a wrapper — Command-as-
Generator, -Transformer, or -Validator — whose only real difference is which
input/output fields the wrapper exposes.  In v1beta1 `Command` is a task
kind like any other, and its data flow is declared on the task itself:

- **`inputs`** — artifact-store paths the command consumes.  The executor
  materializes them in the task temp dir before the command runs and derives
  incoming DAG edges from them.
- **`output`** — the artifact-store path the command produces.  With
  `isStdoutOutput: true` the executor captures stdout and stores it as this
  path; otherwise the command is expected to write the path itself (under
  `buildContext.tempDir`).
- **`command.stdin`** — names one declared input to wire to the command's
  standard input, so filters like `kubectl-slice -f -` need no temp-file
  arguments.
- **`command.args`** — the argument vector.  Args may interpolate
  `buildContext` fields (tempDir, rootDir, leafDir, holosExecutable) to
  locate inputs and outputs, exactly as in v1alpha6.

A command with an `output` generates or transforms; a command with only
`inputs` validates.  The phase is no longer encoded in the schema — it is a
consequence of the task's position in the DAG.

### BuildContext and CUE tags carry over

`BuildContext` carries over from `api/core/v1alpha6/types.go` (lines
108–125) unchanged in spirit: `TempDir`, `RootDir`, `LeafDir`, and
`HolosExecutable`, injected by holos just before evaluating a TaskSet.  The
CUE tag constants also carry over unchanged:
`BuildContextTag = "holos_build_context"`, `ComponentNameTag`,
`ComponentPathTag`, `ComponentLabelsTag`, and `ComponentAnnotationsTag`
(`api/core/v1alpha6/types.go` lines 9–30).  v1beta1 redefines what the tags
inject into (a TaskSet instead of a BuildPlan), not the injection mechanism.

## Design decisions

### D1: Edge derivation

**Decision: DAG edges are derived from `inputs`/`output` matching, and
`dependsOn` adds explicit edges for ordering without data flow.  Both edge
sets are unioned; cycles are an error.**

For every task *T* with input path *p*, the executor adds edge *S* → *T*
for each producing task *S* found by three rules, tried in order:

1. **Exact match** — *S*'s `output` equals *p*.
2. **Directory input** — when no exact match exists, append `/` to *p* and
   prefix-match against all task outputs, carrying the v1alpha6 directory
   rule forward: `inputs: ["out"]` depends on every task producing
   `out/…`.
3. **Directory output** — *p* falls under a produced directory: some *S*
   outputs `d` and *p* begins with `d/`, so a task may consume one file
   from a directory another task produced.

An exact-match input has exactly one producer, because outputs are
write-once — a second task declaring an already-declared output is an error
naming both tasks.  A directory input may legitimately match several
producers and gains an edge from each.  An input matching no output must
exist in the component directory, or edge derivation fails with an error
naming the task and the unmatched path.

`dependsOn` adds edges by task name for constraints the data flow cannot
express: an artifact sink waiting on a validator that produces no output, or
a command that must run after another for side-effect ordering.  `dependsOn`
is struct-keyed rather than list-shaped for the same reason `spec.tasks` is:
two mixins may each unify their own `dependsOn: "validate-x": {}` entry into
the same task without a positional list conflict.  Keys are resolved within
the component's TaskSet; Phase 2 extends resolution to canonical IDs
([D3](#d3-task-naming-and-namespacing)) for cross-component ordering.

A duplicate edge (declared in `dependsOn` and derived from data flow) is
harmless and legal.  Any cycle in the union is an error reported with the
full cycle path.  Rationale: deriving edges from data flow keeps the common
case declaration-free — the pipeline shape *is* the data flow — while
`dependsOn` covers what data flow cannot.  Deriving edges only from
`dependsOn` would force authors to restate every data dependency by hand
and let declarations drift from reality.

### D2: Artifact writing

**Decision: artifact writing is a task kind — `Artifact`, a sink node in
the DAG — not a field on the spec.**

v1alpha6 writes the `artifact:` path implicitly after validators pass; the
gating rule ("validators gate writes") is prose, not structure.  v1beta1
models the write as an explicit task so the gate is a DAG edge like any
other: the producer task outputs a path, a validator `Command` task declares
that path in its `inputs`, and the `Artifact` sink declares the same path in
its `inputs` **and** names the validator in `dependsOn`.  The sink cannot
run until both the producer (data edge) and the validator (explicit edge)
complete successfully; a failed validator fails the render before anything
is written.

Rationale: a sink task keeps the executor uniform — Phase 2's platform-wide
scheduler ([rendering.md](rendering.md)) sees one node type with edges, no
special end-of-pipeline write step — and makes writes composable: a mixin
inserts an additional validator between producer and sink by unifying one
task and one `dependsOn` entry, as the worked example below demonstrates.

### D3: Task naming and namespacing

**Decision: task names are unique within a component's TaskSet; the
platform DAG namespaces tasks by component path with the canonical ID
format `<component-path>:<task-name>`.**

Within one TaskSet, the struct key in `spec.tasks` is the task name — CUE
guarantees uniqueness structurally.  Task names MUST match
`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$` (an RFC 1123 label: lowercase
alphanumerics and interior hyphens — no colons, no slashes), enforced as a
pattern constraint on the `tasks` struct key in the CUE schema and
revalidated by the executor.  Component paths MUST NOT contain a colon,
validated when the platform loads.  When Phase 2 merges all component
TaskSets into one platform DAG, each task's canonical ID is its component's
path relative to the platform root (as injected by `holos_component_path`),
a colon, and the task name: `components/vault:helm`.  The two constraints
above guarantee a canonical ID contains exactly one colon, so it parses
unambiguously and never collides platform-wide (component paths are unique
within a platform).  Bare keys in `dependsOn` resolve within the local
TaskSet; a key containing a colon is a canonical ID referencing another
component's task.  Phase 2 relies on this format for cross-component edges
and log labels; Phase 1 uses only the local names.

### D4: Field naming

**Decision: JSON field names follow the Kubernetes API convention,
lowerCamelCase.**

The v1alpha6 design note's item 1 reads "k8s style lowerSnakeCase", which is
imprecise as written — Kubernetes API conventions name JSON fields in
lowerCamelCase (`apiVersion`, `resourceVersion`); there is no snake_case in
the k8s API surface.  v1beta1 resolves the wording in favor of what the note
plainly meant: the Kubernetes convention.  Every field in the schema above
is lowerCamelCase (`dependsOn`, `isStdoutOutput`, `buildContext`), matching
v1alpha6's existing practice, so no field renames occur on carried-over
structs.  Item 2 of the design note is satisfied structurally:
`spec.tasks` is a struct keyed by name, not a list.

### D5: Open and closed structs

**Decision: `Task` and every kind-specific config struct are closed in CUE;
`spec.tasks` is open (new keys may always be added); a component module's
`#Config` is closed.**

The published CUE schema closes `#TaskSet`, `#TaskSetSpec`, `#Task`,
`#Command`, `#Artifact`, and each carried-over config struct: a typo like
`dependson:` or a field not in the schema fails evaluation with an error
naming the offending field.  Within `#Task`, exactly one kind-specific
config field may be set, and it must match `kind`.  `cue get go` generation
alone does not produce that constraint — it emits independent optional
fields (`helm?: #Helm`, `command?: #Command`) with no union discrimination —
so Phase 1 authors the discriminated union explicitly in the published CUE
schema, guarded per `kind` value (for example, `kind: "Helm"` requires
`helm` and forbids the other config fields).  A task with `kind: "Helm"`
and no `helm` config, or with both `helm` and `command` set, fails
evaluation.  The same per-kind guards enforce the I/O requiredness and
cardinality table in [Task kinds](#task-kinds).  The `tasks` struct stays open in the sense that unification
may always add new task names; that openness is the composition mechanism
(design item 4).  At the author layer, a component module's `#Config` is
`close({...})` as specified in the
[README layer model](README.md#a-component-is-a-cue-module) — closedness is
what makes the site-variance contract enforceable.  Open maps remain only
where the schema is intentionally untyped (`Kustomization`, `Values`,
resource bodies), because typing them would couple holos to kubectl and
chart versions.

## Execution: intra-component now, platform-wide next

Phase 1 (HOL-1493) executes one component's TaskSet at a time: build the
edge set per [D1](#d1-edge-derivation), topologically sort, run tasks
concurrently where the DAG allows, exactly as `holos render component` does
today but over one task type instead of three phase-ordered lists.  Phase 2
([rendering.md](rendering.md)) merges every component's TaskSet into one
platform-wide DAG keyed by canonical ID
([D3](#d3-task-naming-and-namespacing)) and schedules the whole platform
with a high level of concurrency (design-inputs item 5).  Nothing in this
schema is intra-component-only: the merge is a struct union of namespaced
tasks, which is why the schema fixes naming and edges the way it does.

## Worked example: two TaskSets unify into one struct

Struct-keyed tasks are what make TaskSets composable by unification.  The
platform collects every component's TaskSet into one struct keyed by
component path.  Two components and one mixin, in three separate CUE files,
unify without any list-append conflict:

```cue
// components/vault/tasks.cue
package holos

TaskSets: "components/vault": spec: tasks: {
	helm: {
		kind: "Helm"
		helm: chart: {name: "vault", version: "0.28.0", release: "vault"}
		output: "vault.gen.yaml"
	}
	deploy: {
		kind:     "Artifact"
		inputs:   ["vault.gen.yaml"]
		artifact: path: "deploy/vault.yaml"
	}
}
```

```cue
// components/argocd/tasks.cue
package holos

TaskSets: "components/argocd": spec: tasks: {
	helm: {
		kind: "Helm"
		helm: chart: {name: "argo-cd", version: "7.7.0", release: "argocd"}
		output: "argocd.gen.yaml"
	}
	deploy: {
		kind:     "Artifact"
		inputs:   ["argocd.gen.yaml"]
		artifact: path: "deploy/argocd.yaml"
	}
}
```

```cue
// mixins/validate-vault.cue — a policy mixin adds a validator to vault's
// pipeline from a different file, and gates the sink on it.  Impossible
// with v1alpha6 lists; two field unifications with struct-keyed tasks.
package holos

TaskSets: "components/vault": spec: tasks: {
	validate: {
		kind:    "Command"
		inputs:  ["vault.gen.yaml"]
		command: args: ["holos", "cue", "vet", "-f", "vault.gen.yaml"]
	}
	deploy: dependsOn: validate: {}
}
```

The three files unify into one struct.  Both components define tasks named
`helm` and `deploy` — no conflict, because the platform struct namespaces
each TaskSet by component path ([D3](#d3-task-naming-and-namespacing)):
the platform DAG holds `components/vault:helm`, `components/vault:validate`,
`components/vault:deploy`, `components/argocd:helm`, and
`components/argocd:deploy`.  The mixin's `validate` task slots into vault's
DAG through its declared input (a derived edge from `helm`) and gates the
sink by unifying `dependsOn: validate: {}` into `deploy`
([D2](#d2-artifact-writing)) — and because `dependsOn` is itself
struct-keyed, a second mixin could add its own gate to the same sink
without conflicting.  Had these fields been lists, the mixins would have
had to append — exactly the CUE conflict this schema exists to eliminate.
