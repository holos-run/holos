# Modules: CUE Module Packaging and Distribution

This chapter specifies how v1beta1 components are packaged and distributed
as CUE modules: publishing to OCI registries, registry routing, version
selection and platform-layer pinning, offline vendoring, embedded helm
charts, module file access at render time, and britney2-style
reverse-dependency impact checking.  It is the design authority for Phase 3
(HOL-1496, module packaging) and Phase 6 (HOL-1499, the promotion gate);
the numbered design decisions
([M1](#m1-publishing)–[M7](#m7-reverse-dependency-checking)) are cited by
later phases as `modules.md#m1-publishing` and so on.  The layer model and
transparency principle this chapter builds on are fixed in the
[README](README.md#the-layer-model-platform--profile--role--component),
and the TaskSet a component module produces in [schema.md](schema.md).

## Before: the v1alpha6 state being replaced

v1alpha6 has no packaging story.  A component is a directory of CUE files
inside the platform repository; reuse across platforms is copy-paste.  Two
mechanisms this chapter preserves and builds on:

- **The remote-pull chart cache.**  The v1alpha6 Helm generator pulls
  charts into a per-component cache:
  [`internal/component/v1alpha6/v1alpha6.go`](../../../internal/component/v1alpha6/v1alpha6.go)
  (lines 151–160) caches under `vendor/{version}` in the component
  directory and fetches with `helm.PullChart` using the component's
  `Chart.Repository` URL and credentials when the cache path is missing.
  A render on a clean checkout reaches the network, and the chart bytes
  are whatever the repository serves that day — not hermetic.
- **The BuildContext temp-directory wiring.**  `BuildContext`
  ([`api/core/v1alpha6/types.go`](../../../api/core/v1alpha6/types.go)
  lines 108–125) injects `TempDir`, `RootDir`, `LeafDir`, and
  `HolosExecutable` via the `holos_build_context` tag just before
  evaluation, and tasks interpolate the fields into command arguments and
  file paths.  v1beta1 carries this over unchanged
  ([schema.md](schema.md#buildcontext-and-cue-tags-carry-over)); the module
  design adds file resolution *on top of* it, never in place of it.

v1beta1 makes the CUE module the primary packaging format for components:
one module per package, published as an OCI artifact, vendorable to the
filesystem, carrying its charts and other file inputs inside itself so
renders are hermetic and offline-capable.

## The module layout

A holos component package is one CUE module laid out as follows — the
tree is normative, adapting holos-run/holos-paas
`docs/research/distribution-package-repo-layout.md` (summarized in the
[README design inputs](README.md#holos-paas-research-team-repository-layout)):

```
modules/<pkg>/
├── cue.mod/module.cue        # module: "example.com/holos/<pkg>@v0"
├── README.md                 # interface docs: every #Config field
├── config.cue                # closed #Config definition — the package interface
├── component.cue             # the component function: #Config → TaskSet
├── vendor/charts/<chart>-<version>.tgz   # embedded helm chart(s)
└── examples/<name>/          # a reference #Config instance …
    └── deploy/               # … and its committed golden render
```

To be a **valid holos component package**, a module MUST:

1. Declare a domain-qualified module path with a major-version suffix in
   `cue.mod/module.cue` (for example
   `module: "example.com/holos/vault@v0"`).
2. Export a **closed `#Config` definition** — the package's entire
   site-variance interface, per the
   [transparency principle](README.md#normative-the-transparency-principle).
3. Export a **component function** from `#Config` to a `core.#TaskSet` —
   the `#Component` shape fixed in the
   [README layer model](README.md#a-component-is-a-cue-module).
4. Commit every file its tasks consume — chart archives, Kustomize bases,
   CRD schemas — inside the module directory ([M5](#m5-embedded-helm-charts),
   [M6](#m6-module-file-access-at-render-time)).
5. Carry at least one `examples/<name>/` directory holding a reference
   `#Config` instance and its committed golden render — the package's test
   suite and the reverse-dependency gate's input
   ([M7](#m7-reverse-dependency-checking)).
6. Declare no site opinions: no environment names, promotion order,
   compliance mappings, or org-chart roles anywhere in the module.

Degenerate packages are valid: a schema-only module (vendored CRD types)
or a mixin-only module (dashboards, alert thresholds, policy constraint
sets — [resources.md V6](resources.md#v6-unification-by-downstream-modules)).
Requirements 2–3 apply to component packages; the rest apply to every
package kind.

## Design decisions

### M1: Publishing

**Decision: the publish unit is the CUE module, published with
`cue mod publish` to any OCI registry under an immutable semver tag.
Holos adds no publish mechanism of its own.**

`cue mod publish <version>` packages the module directory into a module
zip and pushes it as an OCI image manifest whose config media type,
`application/vnd.cue.module.v1+json`, identifies the artifact as a CUE
module — the manifest carries no `artifactType` field, so tooling that
discriminates CUE modules MUST match the config media type.  Layer 0 is
the module zip; layer 1 is a standalone copy of `cue.mod/module.cue`, so
resolvers walk the dependency graph without downloading module content —
what makes the reverse-dependency graph of
[M7](#m7-reverse-dependency-checking) cheap to compute.

Two publishing semantics matter to holos packages:

- **The zip carries all module files, not only CUE files.**  With
  `source: kind: "git"` in `module.cue` (required at publish time), the
  zip contains the VCS-tracked files under the module directory —
  `cue.mod/` cache subdirectories excluded.  This is how chart archives
  under `vendor/charts/` travel with the module ([M5](#m5-embedded-helm-charts)):
  they are ordinary committed files.  An untracked or gitignored file
  silently stays out of the published zip, so the package lint gate MUST
  verify every file a task references resolves inside the module and is
  tracked by the VCS.
- **Versions are immutable.**  A published `module@version` never changes;
  channels and promotion ([M7](#m7-reverse-dependency-checking)) are
  registry-tag conventions layered above immutable versions.

Consumers resolve module *paths* against a registry and never see the
publishing repository — repository layout stays a private choice of the
publishing team, revisable without breaking consumers.

### M2: Registry routing

**Decision: holos resolves modules with the standard CUE resolver —
`cuelang.org/go/mod/modconfig` semantics via `cuelang.org/go/cue/load` —
honoring `CUE_REGISTRY` and `CUE_CACHE_DIR` exactly as the `cue` CLI does.
Holos MUST NOT implement its own resolver.**

`CUE_REGISTRY` supports comma-separated longest-prefix routing:

```
CUE_REGISTRY='example.com/holos=registry.example.com/cue,registry.cue.works'
```

routes `example.com/holos/...` modules to a private registry and
everything else to the default.  Holos already builds instances through
the standard loader — `BuildInstance` in
[`internal/component/cue.go`](../../../internal/component/cue.go) and
[`internal/cue/cue.go`](../../../internal/cue/cue.go) call
`load.Instances` with a `load.Config` — which delegates registry
resolution to the `modconfig` resolver.  v1beta1 keeps that delegation:
any registry configuration that works for `cue export` works for
`holos render`, auth included.  The one holos-specific registry surface is
the vendor path ([M4](#m4-vendoring)), which substitutes a local directory
for the network without changing resolution semantics.  Rationale: the
resolver is where module identity, security, and caching meet; a bespoke
resolver would fork all three.  The `cue` CLI and holos MUST agree about
which bytes a module path names.

### M3: Version selection and platform pinning

**Decision: version selection is CUE's Minimal Version Selection.  Pins
live in the platform module's `cue.mod/module.cue` `deps` map — the
platform repo is itself a CUE module, and its deps are the site operator's
lockfile-like surface.  The vendor manifest ([M4](#m4-vendoring)) records
the exact resolved build list.**

CUE implements Go-style MVS: `cue.mod/module.cue` declares `deps` as a map
of module paths to concrete minimum versions, and the build list is the
minimum set of versions satisfying all requirements — deterministic, no
lock file.  The platform repository is the main module of every render, so
its `deps` entries are the floor every resolution starts from.  A site
operator pins the vault package like so:

```cue
// <platform-root>/cue.mod/module.cue
module: "platform.example.com@v0"
language: version: "v0.15.1"

deps: {
	"example.com/holos/vault@v0": v: "v0.3.2"
	"example.com/holos/vault-observability@v0": v: "v0.7.1"
}
```

and imports the module from the profile layer as ordinary CUE:

```cue
// profiles/management/vault.cue
package holos

import vault "example.com/holos/vault"

TaskSets: "components/vault": (vault.#Component & {
	config: vault.#Config & {
		namespace: "vault-system"
		version:   "1.17.2"
	}
}).taskSet
```

`cue mod get example.com/holos/vault@v0.3.2` writes the pin; `cue mod tidy`
keeps the deps map consistent with the imports.  Two MVS properties are
stated here so operators are not surprised:

- **A pin is a floor, not a ceiling.**  If another dependency requires
  `vault@v0.4.0`, MVS selects v0.4.0.  The selection is still
  deterministic: a function of the committed deps closure, never of time
  or registry state.
- **The exact build list is recorded at vendor time.**  MVS needs no lock
  file for determinism, but operators want an auditable record; the vendor
  manifest of [M4](#m4-vendoring) is that record, and CI SHOULD fail when
  the manifest disagrees with re-resolution — the same diff-clean
  discipline the golden renders use.

Rationale for pinning at the Platform layer: the platform is one site
instance ([README layer model](README.md#the-layer-model-platform--profile--role--component)),
and version selection is a site decision.  Profiles and roles express
*compatibility* (a major-version import); the platform expresses *choice*
(the pinned minor/patch).  Pins anywhere else would scatter site truth
into reusable layers.

### M4: Vendoring

**Decision: `holos module vendor` materializes the platform's resolved
module dependencies under `vendor/modules/` at the platform root, records
the build list in `vendor/modules/modules.yaml`, and subsequent renders
resolve modules from the vendor tree with no network access.
`holos module pull <module>@<version>` fetches and unpacks one module.**

The on-disk layout mirrors module identity:

```
<platform-root>/vendor/modules/
├── modules.yaml                              # the vendor manifest
└── example.com/holos/
    ├── vault@v0.3.2/                         # one unpacked module zip
    │   ├── cue.mod/module.cue
    │   ├── config.cue
    │   ├── component.cue
    │   └── vendor/charts/vault-0.28.0.tgz
    └── vault-observability@v0.7.1/…
```

- **Source of truth**: `holos module vendor` resolves the platform's full
  dependency closure per [M3](#m3-version-selection-and-platform-pinning),
  fetches each module's OCI artifact through the standard resolver
  ([M2](#m2-registry-routing)), verifies the zip digest against the
  registry manifest, and unpacks it to
  `vendor/modules/<module-path>@<version>/`.
- **The manifest** records, per module: path, version, and the module zip
  digest.  It is the auditable lockfile-like record named in M3.
- **Idempotency**: vendoring is content-addressed.  A directory whose
  manifest entry matches the resolved version and digest is left
  untouched; a mismatch is re-unpacked from scratch; module directories no
  longer in the build list are removed.  A second `holos module vendor`
  run is a no-op and needs no network.
- **Offline renders**: when `vendor/modules/` exists, `holos render`
  resolves module imports from the vendor tree — the loader gets a
  registry view backed by the vendor directory instead of the network —
  and task file access ([M6](#m6-module-file-access-at-render-time)) reads
  module files from the same tree.  A vendored platform renders fully
  offline; CI SHOULD render with the network disabled to keep it honest.
- **`holos module pull`** is the single-module form: fetch, verify, and
  unpack one `<module>@<version>` into the same layout and update the
  manifest.  Vendor is pull over the whole build list plus stale-entry
  removal.

Committing `vendor/modules/` is the site operator's choice: committing
trades repository size for a clone-and-render repo needing no registry;
not committing relies on the manifest digests to make CI vendoring
reproducible.  Both are supported; the manifest is committed either way.

### M5: Embedded helm charts

**Decision: a component module embeds its chart as a `.tgz` archive under
`vendor/charts/` inside the module, and the v1beta1 `Helm` task references
it by module-relative path through a `ModuleFileRef`
([M6](#m6-module-file-access-at-render-time)).  Embedded charts need no
vendor task, no repository config, and no network.**

The chart travels as an ordinary committed file in the module zip
([M1](#m1-publishing)).  No CUE mechanism is involved in moving the bytes.
CUE's `@embed()` attribute — the interpreter holos already enables with
`cuecontext.Interpreter(embed.New())` in
[`internal/component/cue.go`](../../../internal/component/cue.go) — could
load the archive into a CUE value with `type=binary`, but embedding is the
wrong tool: the chart bytes would ride through evaluation and JSON
encoding base64-inflated, only to be written back out to disk for `helm`.
The archive never needs to enter CUE evaluation — tasks read files from
the module directory at render time
([M6](#m6-module-file-access-at-render-time)), and `helm template`
consumes the `.tgz` path directly.

The v1beta1 `Helm` config carries over from v1alpha6
([schema.md](schema.md)) with one addition to `Chart`: a `source` field, a
`ModuleFileRef`, mutually exclusive with `repository`:

```cue
// component.cue, inside the vault module
taskSet: core.#TaskSet & {
	spec: tasks: helm: {
		kind: "Helm"
		helm: chart: {
			name:    "vault"
			version: "0.28.0"
			release: "vault"
			source: {
				module: "example.com/holos/vault"
				path:   "vendor/charts/vault-0.28.0.tgz"
			}
		}
		output: "vault.gen.yaml"
	}
}
```

- When `source` is set, the executor resolves it per
  [M6](#m6-module-file-access-at-render-time) and runs `helm template`
  against the resolved archive path.  No pull occurs and the vendor-task
  deduplication of
  [rendering.md step 3](rendering.md#step-3-deduplicate-helm-chart-vendoring)
  does not apply: the task gets no vendor edge and reads the committed
  chart, exactly as rendering.md already provides for.
- When `repository` is set instead, the remote-pull path carries over:
  the shared vendor task of rendering.md pulls into the platform-level
  cache.  Migration stays incremental, but the packaging convention for
  published modules is the embedded chart — a module whose render reaches
  a chart repository is not hermetic and fails package lint.

Contrast with v1alpha6 (`internal/component/v1alpha6/v1alpha6.go` lines
151–160): the per-component `vendor/{version}` cache plus `helm.PullChart`
put the *cache* in the consumer's tree and got the *bytes* from the
network.  M5 inverts both — the bytes live in the publisher's module,
versioned with it, and a chart upgrade is a reviewable module release
whose golden-render diff shows exactly what changed.

### M6: Module file access at render time

**Decision: tasks reference module files with an explicit `ModuleFileRef`
— `{module, path}` — resolved by the executor to a filesystem location;
resolution enforces containment inside the named module's root.  The
v1alpha6 BuildContext temp-directory wiring carries over unchanged.**

```go
// ModuleFileRef names a file or directory inside a CUE module.
type ModuleFileRef struct {
	// Module is the module path without version suffix, e.g.
	// "example.com/holos/vault".  Empty means the platform's own module.
	Module string `json:"module,omitempty" yaml:"module,omitempty"`
	// Path is the file or directory path relative to the module root.
	Path FileOrDirectoryPath `json:"path" yaml:"path"`
}
```

Resolution, normatively:

1. **Locate the module root.**  An empty `module` resolves to the platform
   root (`buildContext.rootDir`).  A non-empty `module` resolves to the
   vendored directory `vendor/modules/<module>@<version>/`
   ([M4](#m4-vendoring)) — or, unvendored, to the standard CUE module
   cache entry — at the build-list selection of
   [M3](#m3-version-selection-and-platform-pinning).  A module absent from
   the build list is an error naming the task and the module path: a task
   cannot read from a module the platform does not depend on.
2. **Join and contain.**  The executor joins the module root with `path`,
   cleans the result, and verifies containment: `path` MUST be relative,
   and the cleaned join MUST remain inside the module root
   (`filepath.Rel` from the root yields no `..` element).  Symlinks are
   resolved (`filepath.EvalSymlinks`) and the real path re-verified
   against the real module root, so a symlink committed inside a module
   cannot escape it.  A violation fails the render with an error naming
   the task, the module, and the offending path.

Modules are third-party inputs, and the rule's scope is precise:
containment is a guarantee for **declarative file references**
(`ModuleFileRef` resolution), not a sandbox.  A `Command` task executes an
arbitrary program and holos does not confine what that program reads — the
environment, the network, or files outside the module.  The trust boundary
is therefore the `Command` surface: the transparency principle's
no-outside-reads rule is enforced for declarative references by this
algorithm, and for `Command` tasks by lint and review — the package lint
gate MUST surface every `Command` task for review, and a distribution MAY
reject packages whose commands are not on its allowlist.  A render's
declaratively reachable inputs are the platform repo, the vendored
modules, and the values holos injects; `Command` tasks extend that set to
whatever the executed program can reach.

**BuildContext is preserved.**  `TempDir`, `RootDir`, `LeafDir`, and
`HolosExecutable` (`api/core/v1alpha6/types.go` lines 108–125) carry over
with the `holos_build_context` tag injection exactly as
[schema.md](schema.md#buildcontext-and-cue-tags-carry-over) fixes; tasks
keep sharing the per-TaskSet temp directory, and `Command` args keep
interpolating `buildContext` fields.  `ModuleFileRef` is deliberately not
a BuildContext field: BuildContext answers "where is this render
happening" once per TaskSet, while module file resolution answers "where
is this dependency's content" per reference, as a function of the version
selection — so it lives in the executor, keyed by the build list, not in
the injected tag.  Phase 3 extends the carried-over config structs that
accept file paths today (the `File` task source, Kustomize refs) to also
accept a `ModuleFileRef`, without changing the task I/O contract of
[schema.md](schema.md#task-kinds).

### M7: Reverse-dependency checking

**Decision: holos ships britney2-style impact checking as three commands —
`holos module graph`, `holos module rdeps`, and `holos module check` — and
`check` classifies every impact as a schema break (CUE unification
failure, fast) or a behavior break (golden-render diff, slower), reporting
both in machine-readable form for a promotion pipeline.**

Debian's promotion rule — an update that breaks its reverse dependencies
never ships — becomes computationally cheap under the rendered manifest
pattern: renders are hermetic (M5, M6), parallelizable, and need no
cluster.  This decision adapts holos-paas
`docs/research/cue-module-distribution.md` §4.8 into holos commands.
Module dependencies are declared in each module's `cue.mod/module.cue`
`deps` map, exposed as OCI layer 1 ([M1](#m1-publishing)), so graph
construction downloads metadata only.  The **universe** — the set of
modules considered — defaults to the current module and its vendor tree;
a promotion pipeline widens it with `--universe` (a file listing module
paths, or a registry prefix to enumerate).  Within the universe:

- `holos module graph` prints the full dependency graph (one
  `dependent dependency@version` edge per line, or `--format=json`).
- `holos module rdeps <module>` prints the transitive reverse dependents
  of a module — every module whose build list would be affected by a new
  version of it.

**The check workflow.**  `holos module check <module>@<version>
[--baseline <version>]` gates a candidate version:

1. **Enumerate reverse dependents** of `<module>` in the universe, per
   `rdeps`.
2. **Substitute the candidate**: construct each reverse dependent's build
   list with the candidate version replacing the baseline (the currently
   selected version unless `--baseline` names another).
3. **Schema pass (fast, no render)**: re-evaluate each reverse
   dependent's CUE against the candidate — its imports of the candidate's
   definitions and each `examples/<name>/` `#Config` instance.  A
   unification or vet failure is a **schema break**; no render is
   attempted for that dependent.
4. **Behavior pass (slower)**: for dependents that pass, re-render every
   `examples/<name>/` at the candidate version and diff against the
   committed golden `deploy/` tree.  Any diff is a **behavior break** —
   not necessarily wrong, but never silent: an intentional change ships
   with regenerated golden renders, or with a major-version bump.
5. **Report** machine-readably and exit nonzero on any break.

**The report sketch.**  Normative shape, fields finalized in Phase 6
(HOL-1499); the report follows the kind/apiVersion resource convention so
pipelines discriminate it like any other holos output:

```yaml
kind: ModuleCheckReport
apiVersion: v1beta1
candidate:
  module: example.com/holos/vault
  version: v0.4.0
  baseline: v0.3.2
verdict: fail            # pass | fail
counts: {pass: 11, schemaBreak: 1, behaviorBreak: 2, error: 0}
results:
  - module: example.com/holos/vault-observability
    version: v0.7.1
    example: examples/default
    verdict: schemaBreak # pass | schemaBreak | behaviorBreak | error
    detail: |
      #Config.tls: field not allowed
  - module: example.com/profiles/management
    version: v1.2.0
    example: examples/ha
    verdict: behaviorBreak
    diff: {files: 1, insertions: 12, deletions: 4, paths: ["deploy/vault.yaml"]}
```

`error` marks a dependent whose check could not run (fetch failure,
render tooling error) — distinct from a break, and also promotion-
blocking.  A promotion pipeline uses the report directly: publish freely
to `unstable` after package lint; promote to `testing` when the check
passes over the distribution's universe; `stable` adds human release
management.  Channels are registry conventions above immutable versions
([M1](#m1-publishing)); holos owns the check, not the channel bookkeeping.

## Worked example: pin, vendor, render offline

A site operator adopts the vault package at a pinned version:

```console
$ cue mod get example.com/holos/vault@v0.3.2
$ holos module vendor
vendored example.com/holos/vault@v0.3.2 (digest sha256:9f2c…)
vendored example.com/holos/vault-observability@v0.7.1 (digest sha256:41ab…)
wrote vendor/modules/modules.yaml
$ holos render platform
```

The pin lands in the platform's deps map (M3); the profile layer
instantiates the component function with site `#Config` values (the CUE
snippets in M3); vendoring materializes the module — chart archive
included — under `vendor/modules/` (M4); and the render resolves the
module import and the chart's `ModuleFileRef` from the vendor tree with no
network (M5, M6).  When v0.4.0 is published, the operator or promotion
pipeline runs `holos module check example.com/holos/vault@v0.4.0` before
moving the pin, and the report says exactly which profiles break and how
(M7).

## Relation to the other chapters

[schema.md](schema.md) fixes the TaskSet a component function produces and
the BuildContext carry-over M6 preserves; [rendering.md](rendering.md)
step 3 deduplicates remote chart pulls, and M5's embedded charts are the
no-pull case it already anticipates; [resources.md](resources.md) V6 gives
mixin modules their unification surface over rendered output; and
[use-case.md](use-case.md) composes published modules into the canonical
management-cluster profile.
