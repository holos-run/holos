# Self-contained Holos platform design

## Goal

Make a freshly initialized v1beta1 platform render Helm, raw Kubernetes, and
existing Kustomize bases without a copied CUE SDK, a matching external Go
toolchain, or component-local Kustomize wrappers.

## Runtime schemas

`holos init platform v1beta1` seeds a compact local runtime owned by Holos.
It expresses TaskSet assembly directly in ordinary CUE maps and copies only
`cue.mod/module.cue`, rather than the generated `cue.mod/gen` and
`cue.mod/pkg` trees. It deliberately validates Holos task shape while leaving
Kubernetes/CRD schemas as explicit, opt-in imports chosen by component
authors.

## Kustomize bases

The v1beta1 author API gains a declarative Kustomize base input. Its source is
resolved relative to the platform root, not the component root. The renderer
runs Kustomize in that existing directory and writes only the rendered result
to Holos artifacts. `LoadRestrictionsNone` is an explicit generic setting; the
default remains restrictive. Resolved paths and symlinks must remain inside the
platform root.

## Toolchain

Update the repository to Go 1.26. Generating the in-repository CUE API
packages continues to work with that toolchain; initialized platforms do not
need that operation.

## Evidence

Integration tests must prove: a fresh v1beta1 platform contains only its CUE
module declaration; raw Kubernetes and generated Helm values render through
the local runtime; an in-repository Kustomize base outside a component root
renders with the generic base input; traversal outside the platform root is
rejected.
