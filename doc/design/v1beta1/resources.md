# Resources: The Rendered-Resource Round-Trip

This chapter specifies the rendered-resource round-trip: how the manifests
produced by the platform DAG are loaded back into CUE as one typed structure,
and how a platform declares verification policies that unify over that
structure before anything is written to the deploy tree.  It is the design
authority for Phase 4 (HOL-1497); its decisions
([V1](#v1-key-formats)–[V6](#v6-unification-by-downstream-modules)) are cited
by later phases as `resources.md#v1-key-formats` and so on.  It builds on the
[schema.md](schema.md) `Artifact` sink
([D2](schema.md#d2-artifact-writing)) and on the
[rendering.md](rendering.md) platform DAG, whose
[step 6](rendering.md#step-6-the-result) hands this chapter its input.
Every CUE mechanism shown below was verified to evaluate as described with
`cue vet` (cue v0.16.0).

## Before: rendered YAML never re-enters CUE

Today the pipeline is one-way.  `holos render platform` evaluates CUE,
renders manifests, and writes them under the write-to directory — `deploy/`
by default (`internal/holos/constants.go` line 6,
`WriteToDefault = "deploy"`) — and stops.  Nothing loads the rendered output
back into CUE; the final YAML is opaque bytes the moment a generator or
transformer produces it.

The consequence is that no platform-wide property of the *final* output is
checkable.  A v1alpha6 Validator sees one artifact of one component at a
time, so "no `v1/Secret` anywhere on the platform", "every container image
is pinned by digest", or "no two components claim the same final artifact
path with different contents" have no place to be expressed.  The rendered
manifests pattern makes final output reviewable by humans in git diffs; the
round-trip makes it verifiable by machines in CUE — with the same
unification semantics that already govern every other layer of a holos
platform.

## The round-trip structure

After the producer subgraph completes, holos loads every rendered manifest
into one structure keyed by file path, then group/version/kind, then
namespace/name — three levels, so a single resource can be picked out of a
multi-document YAML file:

```cue
// The public API platform engineers unify over.
#Resources: [FilePath=string]: [GVK=string]: [NamespacedName=string]: {...}
```

An instance populated by holos after DAG execution:

```cue
resources: #Resources
resources: "deploy/components/vault/vault.gen.yaml": {
	"apps/v1/Deployment": "vault/vault": {
		apiVersion: "apps/v1"
		kind:       "Deployment"
		metadata: {name: "vault", namespace: "vault"}
		spec: replicas: 1
		// ... the full resource, verbatim
	}
	"v1/Secret": "vault/vault-unseal": {...}
}
```

The closed-ness of the shape is normative.  All three levels of `#Resources`
are pattern constraints over `string`: any key may appear at each level, so
the structure is open in exactly the way `spec.tasks` is open
([schema.md D5](schema.md#d5-open-and-closed-structs)) — openness is the
composition mechanism.  A policy adds constraints by unifying new patterns
or fields into the structure, never by appending to a list.  The innermost
resource body is an open struct (`{...}`), deliberately untyped for the same
reason D5 leaves `Values` and resource bodies open: typing bodies would
couple holos to Kubernetes API versions.  Bodies are loaded verbatim — every
field of the rendered document, unmodified, so what policies see is exactly
what would be applied to a cluster.

`#Resources` is published in the v1beta1 CUE schema package alongside
`#TaskSet`.  No Go struct mirrors it: the structure is data built by the
loader ([V4](#v4-loader-mechanics)), not schema transcribed into
`api/core/v1beta1/types.go`.

### V1: Key formats

**Decision: `FilePath` is the logical artifact path relative to the platform
root; `GVK` is the document's `apiVersion`, a `/`, and its `kind`;
`NamespacedName` is `metadata.namespace`, a `/`, and `metadata.name` when
the namespace is non-empty, and bare `metadata.name` otherwise.  There are
no escaping rules: a key segment containing `/` is a load error.**

- **FilePath** — the `Artifact` sink's final path joined under the default
  write-to prefix (`deploy`), forward-slash separated:
  `deploy/components/vault/vault.gen.yaml`.  The key is *logical*: a
  `--write-to` override relocates bytes on disk but does not change keys, so
  policies that match file paths are stable across invocations and CI
  environments.  Artifact paths are `FileOrDirectoryPath`: a sink whose
  artifact is a directory contributes one FilePath key per manifest file
  under it — the sink's logical path joined with the file's path relative
  to the directory root (`deploy/components/vault/manifests/rbac.yaml`) —
  so every key names a file, never a directory
  ([V4](#v4-loader-mechanics)).
- **GVK** — `apiVersion` is already `group/version` for named groups and
  bare `version` for the core group, so the key is a direct concatenation
  with no discovery or parsing: `apps/v1/Deployment`, `v1/Secret`,
  `gateway.networking.k8s.io/v1/HTTPRoute`.
- **NamespacedName** — `vault/vault` for a namespaced resource, `vault` for
  a cluster-scoped one (a Namespace, a ClusterRole).  No sentinel prefix:
  names cannot contain `/`, so a key containing a slash is unambiguously
  namespaced and a bare key unambiguously carries no namespace.  Scope is
  syntactic, not discovered — holos reads `metadata.namespace` from the
  document and never consults a cluster or an OpenAPI schema, so a
  *namespaced* resource that omits `metadata.namespace` (deferring to
  kubectl's context at apply time) also keys as bare `name`.  Components
  SHOULD set `metadata.namespace` explicitly on namespaced resources;
  policies that must catch both forms match the whole level with `[_]`.

Escaping rules: none exist, by construction.  Kubernetes constrains every
segment — group is a DNS subdomain, version and kind are alphanumeric
identifiers, and namespace and name are path-segment names that reject `/`
— so the only legal slash inside a key component is the one separating
group from version in a non-core `apiVersion`.  The loader validates each
field against exactly that grammar, per field rather than per character:
`apiVersion` MUST contain at most one `/` (a non-empty group before it when
present); `kind`, `metadata.namespace`, and `metadata.name` MUST contain
none.  A violation is an error naming the file path and document index.
The resulting GVK key therefore contains exactly one slash (core group) or
two (named group) and parses unambiguously from the right: the last
segment is the kind, the segment before it is the version, and any
remainder is the group — `apps/v1/Deployment` and `v1/Secret` both parse;
a document with `apiVersion: a/b/c` is rejected at load.  Rather than
define an escape syntax for documents that violate the grammar, the loader
rejects them: an escaping scheme would make keys unpredictable to the
policy author (is it `vault%2Fbad` or `vault\/bad`?) to support only
resources a cluster would refuse anyway.

### V2: Duplicate detection

**Decision: two documents producing the same (FilePath, GVK, NamespacedName)
key is an error, detected structurally by the loader at insertion time —
including byte-identical duplicates that CUE unification would silently
merge.**

Duplicate keys arise within one file: a multi-document YAML artifact
carrying the same resource twice, typically when two generators' outputs
are joined without deduplication.  Duplicates *across* files cannot reach
the loader: final artifact paths are platform-global and write-once, so two
sinks declaring the same path already fail at graph-build time
([rendering.md step 2](rendering.md#step-2-collect-tasksets-and-merge-into-one-dag)).
This chapter extends that rule to prefix-freedom: no sink's final path may
equal *or nest under* another sink's directory path — a file sink at
`components/vault/rbac.yaml` conflicts with a directory sink at
`components/vault` — validated at graph build with an error naming both
canonical IDs.
Prefix-freedom guarantees every FilePath key has exactly one producing
sink and protects the atomic directory promotion of
[rendering.md R8](rendering.md#r8-failure-semantics), which could not
swap a directory another sink writes into.
The same resource appearing in two *different* files is not a duplicate —
the file-path level exists precisely so both load — though a policy may
forbid it by unifying a constraint over the structure.

The error message contract names the file path, the GVK, the
namespace/name, and both document positions:

```
error: duplicate resource: deploy/components/vault/vault.gen.yaml:
apps/v1/Deployment vault/vault: document 4 duplicates document 1
```

Why a structural check rather than bare unification: unification catches
*conflicting* duplicates for free — two different bodies at one key are a
CUE conflict naming both values — but it silently merges identical bodies,
and a resource emitted twice is a pipeline bug worth surfacing even when
the copies agree.  The loader builds the nested structure in Go maps before
encoding into CUE ([V4](#v4-loader-mechanics)), so the insertion check is
where duplicate keys surface with document indexes still in hand.  CUE
conflict detection remains in force above the loader: when a *policy*
constrains a key, the conflict error is the verification mechanism working
as intended ([V3](#v3-verification-policy)).

### V3: Verification policy

**Decision: a platform declares policy packages in `spec.policies`,
struct-keyed by name.  Holos verifies resources at a synthetic barrier node
between producers and sinks, so a failing policy fails the render with a
non-zero exit after all artifacts are built but before anything is written
to the deploy tree.**

A policy is an ordinary CUE package within the platform's module.  The
platform declares each one by path, keyed by name (structs, not lists —
design-inputs item 2 — so mixins may add policies by unification):

```cue
// Platform spec (v1beta1).
spec: policies: "no-secrets": path: "./policy/no-secrets"
```

The field is normative for the v1beta1 `PlatformSpec`; Phase 4 (HOL-1497)
transcribes it into `api/core/v1beta1/types.go` alongside the fields the
spec carries over from `api/core/v1alpha6/types.go`:

```go
// PlatformSpec gains one field in v1beta1.
type PlatformSpec struct {
	// ... existing fields carry over from v1alpha6.

	// Policies represents verification policy packages keyed by policy
	// name, unified over the rendered resources before artifacts are
	// written (resources.md V3).
	Policies map[string]Policy `json:"policies,omitempty" yaml:"policies,omitempty"`
}

// Policy represents one verification policy CUE package.
type Policy struct {
	// Path represents the policy package directory relative to the
	// platform root, e.g. "./policy/no-secrets".
	Path string `json:"path" yaml:"path"`
}
```

Policy names MUST match the same RFC 1123 label constraint task names
follow ([schema.md D3](schema.md#d3-task-naming-and-namespacing)), so a
failing policy's name composes into error messages and log labels
unambiguously.  A policy path MUST be relative, forward-slash separated,
and resolve inside the platform's CUE module: absolute paths and paths
whose normalized form escapes the platform root (`..`) are rejected when
the platform loads — the same containment rule component paths follow.
Each package is built with `BuildInstance(root, path)` exactly as a
component is, so a policy package may import any module the platform's
`cue.mod` resolves ([V6](#v6-unification-by-downstream-modules)).

The package declares a `resources` field.  After the loader builds the
populated structure, holos unifies each policy package's `resources` value
with it and validates the result; any conflict fails the render.  The
worked example [below](#worked-example-forbid-v1secret-platform-wide) shows
the mechanism end to end.

**Timing.**  Written-then-failed is rejected.  Verification runs at a
synthetic barrier node, `holos.internal:verify-resources` (under the
reserved virtual component path of
[rendering.md R6](rendering.md#r6-shared-vendor-task-identity)), with an
incoming edge from every producer of every `Artifact` sink's input and an
outgoing edge to every `Artifact` sink.  This is
[D2](schema.md#d2-artifact-writing)'s validators-gate-writes rule promoted
platform-wide: when a policy fails, the failure cancels the render per
[R8](rendering.md#r8-failure-semantics), no sink has run, the deploy tree
is untouched, and the previous complete render is preserved.  The cost is a
global synchronization point — no sink writes until the slowest producer
platform-wide completes.  Sinks are cheap atomic writes, so the delay is
negligible, and the benefit is categorical: a policy-violating resource
never reaches the deploy tree, so a git checkout of a rendered platform
that passed CI never contains a violation.

**Mechanism.**  A policy forbids a value by constraining a field to a
conflicting string literal, following the repository's prior art
(`cmd/holos/tests/cli/cue-vet.txt`,
`secret: kind: "Forbidden. Use an ExternalSecret instead."`).  An explicit
`_|_` also fails, but was rejected because its error —
`explicit error (_|_ literal) in source` — names only the policy's source
line and drops the resource path entirely, while a literal conflict names
the full CUE path (file, GVK, namespace/name), both values, and both
positions.  Optional-field constraints (`"v1/Secret"?:`) scope a policy to
resources that are present without requiring any to exist.

Verification is a `holos render platform` concern.  `holos render
component` renders one TaskSet without platform context and applies no
platform policies; component-scoped validation remains the province of
validator `Command` tasks
([schema.md Task kinds](schema.md#task-kinds)).

### V4: Loader mechanics

**Decision: the parent process loads every sink input from the artifact
store once, after all producers complete; documents are split with a
`yaml.v3` decode loop; the structure is built in Go maps, encoded into CUE,
and unified with each policy package built by the existing `internal/cue`
loader.**

- **Source** — the loader reads each `Artifact` sink's declared input from
  the artifact store, not from files on disk: at the barrier node the deploy
  tree has deliberately not been written ([V3](#v3-verification-policy)),
  and the store bytes are exactly what the sinks will write.
- **Directory artifacts** — a sink whose artifact is a directory
  (`FileOrDirectoryPath`) is walked in the store: each file whose name ends
  in `.yaml` or `.yml` loads under its own FilePath key, the sink's logical
  path joined with the file's directory-relative path
  ([V1](#v1-key-formats)); files with any other extension are ancillary
  outputs, skipped and logged at debug level, never resources.  The
  prefix-freedom rule ([V2](#v2-duplicate-detection)) guarantees the walk
  of one sink cannot produce a key another sink also produces.
- **Splitting** — each manifest file is split into documents with a
  `gopkg.in/yaml.v3` `Decoder` loop, the multi-document behavior the
  repository already depends on.  Empty and null documents are skipped.  A
  document that is not a map, or that lacks `apiVersion`, `kind`, or
  `metadata.name`, is an error naming the file path and document index.
- **List flattening** — a document with `apiVersion: v1` and `kind: List`
  is a wrapper `kubectl apply -f` accepts, not a resource: the loader
  flattens it, treating each element of `items` exactly like a top-level
  document — same required fields, same key derivation, same duplicate
  detection ([V2](#v2-duplicate-detection)) — attributed as `document N
  item M` in errors.  The wrapper itself never appears in the structure,
  so a List cannot smuggle a resource past a policy.  Flattening is one
  level: a `v1/List` nested inside another is an error.  Typed collections
  (`kind: DeploymentList` and friends) are not flattened; they fail the
  `metadata.name` requirement, and the error message for a kind ending in
  `List` that carries `items` says to emit the items as separate documents
  or a `v1/List`.  Rendered manifests SHOULD NOT emit wrappers at all —
  multi-document YAML is the pattern's native shape — but a chart that
  does must not break the round-trip, and must not bypass it either.
- **CUE loading** — policy packages are built with `BuildInstance`
  (`internal/cue/cue.go`), the same loader every other holos evaluation
  uses.  `BuildInstance` serializes on `cueMutex` (`cue.go` lines 26–27)
  because CUE evaluation is not safe for concurrent use — and that is why
  the load happens once, in the parent process, after the producer subgraph
  completes: one evaluation, no contention, no compiler subprocess needed.
- **Memory** — the artifact store already holds every rendered artifact in
  memory (`MapStore`, `internal/artifact/artifact.go`); the decoded Go maps
  and the CUE value roughly triple resident manifest bytes during
  verification.  A large platform — say 5,000 resources averaging 4 KiB,
  20 MiB of YAML — costs on the order of 100 MiB transiently, well inside
  the memory envelope CUE evaluation already sets
  (`internal/platform/platform.go` limits concurrency "due to CUE memory
  usage concerns").  Partitioning the load per file would cap the footprint
  but forecloses cross-file policies, which are the point of the round-trip,
  so the whole structure loads in one evaluation.

### V5: CLI surface

**Decision: `holos show resources` renders the producer subgraph and prints
the populated structure in yaml or json, taking the same selector flags as
`holos show buildplans`.**

The command joins `platform` and `buildplans` under `holos show`
(`internal/cli/show.go`), with the flag conventions of `showBuildPlans`:
`--format yaml|json` (default yaml), the component label selector flags
from the shared platform flag set, and `--concurrency`.  It executes the
producer subgraph — charts render, commands run — but stops at the barrier:
no policies are applied and nothing is written.  The output is exactly the
structure policies unify over, which makes the command the policy author's
development loop:

```console
$ holos show resources --format yaml > resources.yaml
$ cue vet ./policy/no-secrets resources.yaml
```

iterates on a policy against a captured render without re-rendering, and
`holos show resources --selector app.holos.run/name=vault` scopes the view
to one component's files while debugging a violation.

### V6: Unification by downstream modules

**Decision: policies compose as ordinary CUE values.  A downstream module
publishes constraint sets as mixin definitions over `#Resources`; a
site-owned policy package imports and unifies them, and the platform
declares only the site-owned package.**

A security team's module exports a constraint set — the mixin-definition
export of the [README layer model](README.md#design-inputs):

```cue
package security

import core "github.com/holos-run/holos/api/core/v1beta1"

// #NoSecrets forbids v1/Secret in every rendered file.
#NoSecrets: core.#Resources
#NoSecrets: [_]: "v1/Secret"?: [_]: kind: "Forbidden: use ExternalSecret instead"
```

The platform's policy package adopts it, and may unify several such
mixins into one `resources` value:

```cue
package policy

import "example.com/security@v1"

resources: security.#NoSecrets
```

Verified with `cue vet`: the three sources — the mixin definition, the
site adoption, and the loaded data — unify into the same conflict error a
locally-authored policy produces, with all three positions cited.  This
composition is the hook the later phases build on: Phase 5 role-layer
transformations receive the same structure through the same unification
seam, and the three-module composition (platform, distribution, security)
demonstrated in [use-case.md](use-case.md) is this decision exercised
across module boundaries.  The distribution of such modules is
[modules.md](modules.md)'s subject.

## Worked example: forbid v1/Secret platform-wide

The policy package, declared as `spec: policies: "no-secrets": path:
"./policy/no-secrets"`:

```cue
// policy/no-secrets/policy.cue
package policy

// Unifies over every rendered file.  Any v1/Secret's kind field conflicts
// with the literal below, failing the render with an error that names the
// file, GVK, and namespace/name.
resources: [_]: "v1/Secret"?: [_]: kind: "Forbidden: use ExternalSecret instead"
```

A component renders a chart that emits a Secret.  After the producer
subgraph completes, the loader populates the structure (shown as CUE for
illustration; holos builds it in memory):

```cue
resources: "deploy/components/vault/vault.gen.yaml": {
	"apps/v1/Deployment": "vault/vault": {
		apiVersion: "apps/v1"
		kind:       "Deployment"
		metadata: {name: "vault", namespace: "vault"}
		spec: replicas: 1
	}
	"v1/Secret": "vault/vault-unseal": {
		apiVersion: "v1"
		kind:       "Secret"
		metadata: {name: "vault-unseal", namespace: "vault"}
	}
}
```

Unifying the policy over the structure — reproducible today with
`cue vet ./policy/no-secrets resources.cue` against the data above — yields
the conflict, and the render fails before any sink runs:

```console
$ holos render platform
error: policy no-secrets: resources."deploy/components/vault/vault.gen.yaml"."v1/Secret"."vault/vault-unseal".kind: conflicting values "Forbidden: use ExternalSecret instead" and "Secret":
    ./policy/no-secrets/policy.cue:7:42
exit status 1
```

Every element of the error contract is present: the file path
(`deploy/components/vault/vault.gen.yaml`), the GVK (`v1/Secret`), the
namespace/name (`vault/vault-unseal`), both conflicting values, and the
policy source position — all carried by the CUE conflict error itself,
because the key formats of [V1](#v1-key-formats) put the address *in the
path*.  The policy's source position is cited; the resource side carries no
file position because the structure is built from memory, which is why the
error message is prefixed with the policy name.  The compliant remainder of
the platform renders unchanged once the Secret is replaced: the same
policy, vetted against the structure with the Secret removed, evaluates
cleanly (verified with `cue vet`, exit 0).

## A public API with compatibility expectations

The `#Resources` structure is a public API for platform engineers, on the
same tier as the TaskSet schema: policies, mixin modules, and downstream
tooling are written against it, and they must keep evaluating for the life
of v1beta1.  Concretely:

- The three-level shape and the key formats of [V1](#v1-key-formats) are
  frozen for v1beta1.  No level is inserted, removed, or reordered; key
  formats do not change.
- Resource bodies remain verbatim and open.  Holos never rewrites, defaults,
  or prunes fields of a loaded resource; a policy written against any body
  field keeps seeing what the manifest says.
- Extensions are additive and arrive only in new key spaces (for example, a
  future metadata level under a reserved key) with a design revision to
  this chapter; a policy that matches with `[_]` patterns and optional
  fields is forward-compatible by construction.
- The duplicate-detection ([V2](#v2-duplicate-detection)) and
  fail-before-write ([V3](#v3-verification-policy)) guarantees are part of
  the contract: tooling may assume a deploy tree produced by a successful
  `holos render platform` satisfied every declared policy.

A breaking change to any of these requires a new API version, exactly as it
would for the core schema.
