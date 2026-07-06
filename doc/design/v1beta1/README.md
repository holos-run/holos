# Holos v1beta1 Design

This document set specifies the v1beta1 API and the architecture that
supports it.  v1beta1 completes the redesign started in v1alpha6: it replaces
the BuildPlan with a composable TaskSet, composes every component's TaskSet
into one platform-wide DAG executed concurrently, and packages components as
independently published CUE modules.

This README owns three things: the document map, the design-inputs digest,
and the layer model chapter.  Every other chapter lives in a sibling file.
Any chapter that grows past roughly 400 lines must be split into additional
sibling files and added to the document map.

## Document map

| Document | Chapter | Status | Owning task |
| -- | -- | -- | -- |
| [README.md](README.md) | Document map, design inputs, layer model | complete | HOL-1501 |
| [schema.md](schema.md) | TaskSet, Task, and first-class Command core schema | stub | HOL-1502 |
| [rendering.md](rendering.md) | Compiler pool and the platform-wide DAG | stub | HOL-1503 |
| [resources.md](resources.md) | Rendered-resource round-trip | stub | HOL-1504 |
| [modules.md](modules.md) | CUE module packaging and distribution | stub | HOL-1505 |
| [use-case.md](use-case.md) | Canonical management-cluster use case | stub | HOL-1506 |

Design review and the phase exit gate are owned by HOL-1507.

## Design inputs

Three inputs seed this design: the in-repo v1alpha6 design note, and two
research documents from the holos-run/holos-paas repository.  The external
documents are cited and summarized here — consult the originals for the full
argument.

### The v1alpha6 design note

[`tasks/v1alpha6-design.md`](../../../tasks/v1alpha6-design.md) recorded the
lessons of v1alpha5 and set out a five-item plan, quoted verbatim:

```markdown
1. Standardize on k8s style lowerSnakeCase for field names.
2. Replace lists with structs, e.g. Platform.spec.componets.
3. Deprecate BuildPlan, use a TaskSet instead.
4. Ensure TaskSets are composable into one big TaskSet for all platform components.
5. Execute the tasks in the big TaskSet using topological sort over the DAG with a high level of concurrency.
```

v1alpha6 shipped items 1 and 2 only partially: its `BuildPlanSpec.Artifacts`
field in [`api/core/v1alpha6/types.go`](../../../api/core/v1alpha6/types.go)
is still a list (`spec.artifacts`), and its field naming remains camelCase
(`apiVersion`, `buildContext`, `tempDir`).  Items 3, 4, and 5 were deferred
entirely — v1alpha6 still renders one BuildPlan per component with no
cross-component task graph.

**v1beta1 implements items 3–5**: BuildPlan is deprecated in favor of the
TaskSet ([schema.md](schema.md)), TaskSets compose into one platform-wide
TaskSet, and that composed TaskSet executes by topological sort over the DAG
with a high level of concurrency ([rendering.md](rendering.md)).  Where
v1beta1 touches fields that v1alpha6 left unfinished, it also finishes items
1 and 2 for the new schemas.

### holos-paas research: CUE Modules as a Package Ecosystem

holos-run/holos-paas `docs/research/cue-module-distribution.md` surveys prior
packaging ecosystems (Puppet, Chef, Ansible, Debian, Helm) and the current
CUE module system, then makes design recommendations.  The points v1beta1
adopts:

- **The package unit is a CUE module** (its §4.1), published as an OCI
  artifact via `cue mod publish`, exporting up to three things: a **closed
  `#Config` definition** (the package's typed interface), a **component** (a
  function from `#Config` values to rendered output — a TaskSet in v1beta1
  terms), and **optional mixin definitions** that other packages unify
  against (dashboards, alert thresholds, policy constraint sets).  Degenerate
  forms — schema-only or policy-only packages — are valid.
- **The transparency principle** (its §4.2): the interface between a package
  and the platform is small, structural, and silent about policy — the
  `io.Reader` of configuration.  A component is a pure function: CUE values
  in, manifests out.  It reads nothing outside its module; all variance
  arrives through `#Config` fields or `@tag()` injection.  Packages never
  define enterprise concepts — no environment names, promotion order,
  compliance mappings, or org-chart roles.  Consumers unify additional
  constraints over a component's output as the sanctioned, type-checked
  escape hatch.
- **The three-layer composition model** (its §4.3), translating Puppet's
  roles-and-profiles pattern: distribution packages (no site opinions) are
  consumed only through site-owned profile modules, which the platform spec
  references.  Environments are a profile-layer concept; site truth lives in
  the consumer's repository, structurally protected because CUE unification
  has no overwrite.  The layer model chapter below adapts this table to
  Holos-native terms.
- **The britney2-style promotion gate** (its §4.8): packages publish freely
  to an `unstable` channel after a lint gate; promotion to `testing` requires
  re-rendering every reverse dependency — every profile and reference
  consumer that depends on the package — at the candidate version and
  requiring those renders and validators to pass.  Debian's "an update that
  breaks its reverse dependencies never ships" promise becomes
  computationally trivial under the rendered manifest pattern, because
  renders are hermetic, parallelizable, and need no cluster.

### holos-paas research: Team Repository Layout

holos-run/holos-paas `docs/research/distribution-package-repo-layout.md`
turns those recommendations into a concrete repository layout for a
contributing team.  The points v1beta1 adopts:

- **One repo per team, one CUE module per package under `modules/<pkg>/`.**
  Each module directory carries its own `cue.mod/module.cue`, a `README.md`
  documenting every `#Config` field, `config.cue` (the closed `#Config`),
  `component.cue` (the component function), and vendored upstream inputs
  (Helm charts, CRD schemas) committed inside the module so renders are
  hermetic.  Component, mixin, and schema packages are sibling modules, never
  subpackages, so consumers can adopt each independently.
- **Committed golden renders are the test suite.**  Each package ships
  `examples/<name>/` directories holding a reference `#Config` instance and
  its committed rendered output; CI re-renders every example and fails on any
  diff.  The examples triple as unit tests, documentation, and
  promotion-gate input for the britney2-style pipeline.
- **CI is stamped, not hand-rolled**: lint, compat-gate, test, and publish
  workflows come from a shared distribution template repository
  (modulesync-style), with per-module path-prefixed release tags
  (`modules/valkey/v1.4.0`).
- **What stays out of a package repo**: environment names, `#Config` values
  for real deployments, tag-injected site data, secrets, rendered deploy
  trees, and admission bindings — all profile- or platform-layer concerns.

[modules.md](modules.md) develops the packaging design; the layer model
below fixes the vocabulary.

## The layer model: Platform → Profile → Role → Component

v1beta1 names four layers.  Two of them exist in Go; two of them exist only
in CUE.

| Layer | Meaning | Authored by | Go tooling awareness |
| -- | -- | -- | -- |
| Platform | one site instance | site operators | yes — `internal/platform/` |
| Profile | cluster class (management, nonprod workload, prod workload) | site operators | **none — pure CUE convention** |
| Role | reusable capability composed of components (the Puppet "profile" pattern, correctly named) | distribution / platform teams | **none — pure CUE convention** |
| Component | one deployable unit; a CUE module | module authors | yes — `internal/component/` |

A Platform is one site instance: the single entrypoint a site operator
renders.  A Profile is a cluster class — a management cluster, a nonprod
workload cluster, a prod workload cluster — selecting which roles apply to
clusters of that class.  A Role is a reusable capability composed of
components: "ingress", "observability", "postgres with backups and
dashboards".  This is the Puppet "profile" pattern with its correct name —
Puppet called the capability layer a profile and the box-level layer a role;
Holos assigns the names the other way so that Role means the reusable
capability.  A Component is one deployable unit, packaged as a CUE module.

### Normative: the Go tooling knows only Platform → Component

**The Go tooling knows only Platform → Component.  Profile and Role are pure
author-layer CUE conventions and MUST NOT be hard-coded in Go.**  No Go
type, no CLI flag, no rendering behavior may depend on the existence,
naming, or structure of profiles or roles.  A site that organizes its CUE
with different intermediate layers — or none — renders identically.
Profiles and roles exist so that humans can compose and review
configuration; by the time the Go tooling evaluates CUE, they have already
collapsed into a Platform holding components.

The current code establishes the two Go-visible layers precisely:

- `internal/platform/platform.go` — `Platform.Load` discriminates the API
  version of the platform resource and dispatches: it reads the type meta,
  switches on `APIVersion` (selecting `v1alpha6.Platform`, defaulting to
  `v1alpha5.Platform`), builds the CUE instance, and loads the `holos` field
  value.  v1beta1 adds one case to this dispatch.
- `internal/component/component.go` — `(*Component).TypeMeta()` reads the
  component's `typemeta.yaml` discriminator file and returns a v1alpha5
  BuildPlan TypeMeta when the file does not exist.  Per-component version
  dispatch keys off this value, so components of different API versions
  coexist in one platform during migration.
- `internal/holos/constants.go` — `TypeMetaFile` names the discriminator
  file (`typemeta.yaml`).

Nothing in `internal/` refers to profiles or roles today, and this design
keeps it that way.

### A Component is a CUE module

A Component is a CUE module exporting a closed `#Config` definition and a
function from `#Config` to a `TaskSet`:

```cue
package component

// Closed site-variance interface.  ALL variance flows through #Config or @tag().
#Config: close({
	namespace: string | *"example"
	version:   string | *"1.2.3"
})

// The component function: #Config in, TaskSet out.
#Component: {
	config: #Config
	taskSet: core.#TaskSet & {
		metadata: name: "example"
		spec: tasks: {...}
	}
}
```

The closed `#Config` is the component's entire site-variance interface:
every field a site may set, with types, constraints, and defaults.
Closedness makes the contract enforceable — a profile that sets a field the
component does not declare fails the render with a CUE error naming both
sources.  The `TaskSet` the function produces is the v1beta1 unit of work
([schema.md](schema.md)); the platform composes every component's TaskSet
into one DAG ([rendering.md](rendering.md)).

### Normative: the transparency principle

**Components are pure functions of `#Config`; modules carry no site
opinions.**  A component MUST NOT read the cluster, the environment, or the
filesystem outside its own module — everything it consumes (vendored charts,
Kustomize bases, CRD schemas) is committed inside the module, and all
variance arrives through `#Config` or `@tag()` injection.  A component MUST
NOT declare environment names, promotion order, compliance mappings, or
org-chart roles; those belong to the profile and platform layers.  This is
what makes renders hermetic and cacheable, and what lets two organizations
with incompatible policies consume the same module unchanged: each unifies
its own constraints over the component's output in its own profile layer.

### How the layers compose

Site operators author the Platform and its Profiles; distribution and
platform teams author Roles; module authors publish Components.  A profile
selects roles for a cluster class; a role selects components and sets their
`#Config` values; the platform instantiates profiles per cluster.  Every
site opinion is greppable in the profile and platform layers because the
component layer structurally has no place to hold one.  The canonical
worked example — a management-cluster profile — is the subject of
[use-case.md](use-case.md).
