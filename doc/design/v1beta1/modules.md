# Modules: CUE Module Packaging and Distribution

> Status: stub — owned by HOL-1505.

This chapter specifies how v1beta1 components are packaged and distributed
as CUE modules: one module per package exporting a closed `#Config`, a
component function producing a TaskSet, and optional mixin definitions;
published as OCI artifacts via `cue mod publish` with `CUE_REGISTRY`
prefix routing; laid out in team repositories under `modules/<pkg>/` with
vendored upstream inputs and committed golden renders as the test suite.
It also covers reverse-dependency render checks — the britney2-style
promotion gate — and the compatibility gate for exported definitions across
versions, summarized in the design-inputs section of the
[README](README.md).
