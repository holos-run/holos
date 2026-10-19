# Self-contained Holos platform implementation plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Remove mandatory generated CUE vendoring from v1beta1 platform initialization and support existing platform-root Kustomize bases.

**Architecture:** Generated platforms use a compact local runtime CUE schema for Holos TaskSet construction. The renderer gains a validated base-directory Kustomize task that preserves source-relative resources while writing only artifacts. Kubernetes/CRD schemas remain optional component dependencies.

**Tech Stack:** Go 1.26, CUE, Holos v1beta1, Kustomize, Helm.

---

### Task 1: Upgrade the supported Go toolchain

Update `go.mod` to Go 1.26, regenerate the in-repository CUE API packages, then run focused CUE and CLI tests.

### Task 2: Make v1beta1 init self-contained

Write an integration fixture that initializes and renders raw Kubernetes and Helm components with no `cue.mod/pkg` or `cue.mod/gen`. Add the compact local runtime template and change init generation until the fixtures write non-empty artifacts.

### Task 3: Add platform-root Kustomize bases

Write failing integration tests for an existing base outside the component root and for rejected traversal outside the platform root. Add the author field, root containment validation, and Kustomize execution path.

### Task 4: Verify and document

Run formatting, focused suites, full tests, and a fresh external-dns plus mediastorage reproduction. Update docs and commit only with fresh evidence.
