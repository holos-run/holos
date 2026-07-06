# Rendering Migration: v1alpha5 and v1alpha6 Components

This chapter specifies how platforms holding v1alpha5 and v1alpha6
components render under the v1beta1 platform-wide DAG of
[rendering.md](rendering.md), so migration proceeds
component-by-component with no flag day.

Per-component API version dispatch is unchanged: `(*Component).TypeMeta()`
reads each component's `typemeta.yaml` discriminator
(`internal/component/component.go`, cited in the
[README](README.md#normative-the-go-tooling-knows-only-platform--component)),
so components of different API versions coexist in one platform.  Each
v1alpha5/v1alpha6 component becomes exactly one opaque node in the platform
DAG — "render component `<path>`", executed through the existing JSON
compiler protocol and the phase-ordered `Build` quoted in
[rendering.md](rendering.md#before-the-v1alpha6-render-pipeline) — with no
intra-component task visibility, no cross-component edges, and no vendor
dedup.  v1beta1 components join the graph natively; opaque nodes schedule
under the same executor and the same
[R7](rendering.md#r7-deterministic-ordering)/[R8](rendering.md#r8-failure-semantics)
rules, so migrating a component changes its graph granularity and nothing
else.  The JSON protocol and the per-component `vendor/{version}` cache
retire together with v1alpha6 support, not before.
